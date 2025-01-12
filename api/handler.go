package api

import (
	"archive/zip"
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
	"bytes"
	"log"
)

// PriceData используется и при чтении CSV, и при выгрузке данных из БД.
type PriceData struct {
	ID        string
	CreatedAt time.Time
	Name      string
	Category  string
	Price     int64 // Цена хранится в копейках для удобства расчетов.
}

// PostResponse описывает JSON-ответ для POST-запроса.
type PostResponse struct {
	TotalItems      int     `json:"total_items"`      // Кол-во реально вставленных строк
	TotalCategories int     `json:"total_categories"` // Общее кол-во категорий в БД
	TotalPrice      float64 `json:"total_price"`      // Суммарная стоимость в формате float
}

// Обработчик для POST-запроса
func UploadPricesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		// Обработка загружаемого архива
		file, _, err := r.FormFile("file")
        if err != nil {
            log.Printf("Error retrieving file: %v", err)
            http.Error(w, "Failed to retrieve file", http.StatusBadRequest)
            return
        }
        defer file.Close()

        buf := make([]byte, 512) // Считываем первые 512 байт для определения типа
        _, err = file.Read(buf)
        if err != nil {
            log.Printf("Error reading file header: %v", err)
            http.Error(w, "Failed to read file header", http.StatusInternalServerError)
            return
        }

        // Считываем ZIP-архив в буфер
        tempFiles := &bytes.Buffer{}
        if _, err := io.Copy(tempFiles, file); err != nil {
            log.Printf("Error reading file: %v", err)
            http.Error(w, "Failed to read file", http.StatusInternalServerError)
            return
        }

        // Открываем ZIP-архив из буфера
        zipReader, err := zip.NewReader(bytes.NewReader(tempFiles.Bytes()), int64(tempFiles.Len()))
        if err != nil {
            log.Printf("Error opening zip: %v", err)
            http.Error(w, "Invalid zip file", http.StatusBadRequest)
            return
        }

		var totalItems int
		var totalCategories int
		var totalPrice int64
		categorySet := make(map[string]bool)

		// Проходим по всем файлам в архиве
		for _, file := range zipReader.File {
			f, err := file.Open()
			if err != nil {
				http.Error(w, "Failed to open file in archive", http.StatusInternalServerError)
				return
			}
			defer f.Close()

			reader := csv.NewReader(f)
			// Пропускаем заголовок
			if _, err := reader.Read(); err != nil {
				http.Error(w, "Failed to read CSV header", http.StatusInternalServerError)
				return
			}

			// Читаем данные построчно
			for {
				record, err := reader.Read()
				if err == io.EOF {
					break
				}
				if err != nil {
					http.Error(w, "Failed to read CSV", http.StatusInternalServerError)
					return
				}

				// Обработка данных из строки
				productID, err := strconv.Atoi(record[0])
				if err != nil {
					http.Error(w, "Invalid product ID", http.StatusBadRequest)
					return
				}

				name := record[1]
				category := record[2]

				priceFloat, err := strconv.ParseFloat(record[3], 64)
				if err != nil {
					http.Error(w, "Invalid price format", http.StatusBadRequest)
					return
				}
				price := int64(priceFloat * 100) // Конвертируем в копейки.

				createdAt, err := time.Parse("2006-01-02", record[4])
				if err != nil {
					http.Error(w, "Invalid date format", http.StatusBadRequest)
					return
				}

				// Записываем в базу
				_, err = db.Exec("INSERT INTO prices (id, created_at, name, category, price) VALUES ($1, $2, $3, $4, $5)",
					productID, createdAt, name, category, price)
				if err != nil {
					http.Error(w, "Failed to insert data into database", http.StatusInternalServerError)
					return
				}

				// Увеличиваем счетчики
				totalItems++
				totalPrice += price
				categorySet[category] = true
			}
		}

		totalCategories = len(categorySet)

		// Возвращаем JSON
		response := PostResponse{
			TotalItems:      totalItems,
			TotalCategories: totalCategories,
			TotalPrice:      float64(totalPrice) / 100, // Конвертируем обратно в формат float.
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// Обработчик для GET-запроса
func DownloadPricesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		// Создаем временный CSV-файл
		csvFile, err := os.CreateTemp("", "data-*.csv")
		if err != nil {
			http.Error(w, "Failed to create CSV file", http.StatusInternalServerError)
			return
		}
		defer os.Remove(csvFile.Name())

		writer := csv.NewWriter(csvFile)
		defer writer.Flush()

		// Запрос данных из базы
		rows, err := db.Query("SELECT id, created_at, name, category, price FROM prices")
		if err != nil {
			http.Error(w, "Failed to query database", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// Записываем заголовки
		writer.Write([]string{"id", "name", "category", "price", "create_date"})

		// Записываем строки
		for rows.Next() {
			var productID int
			var createdAt time.Time
			var name, category string
			var price int64

			err := rows.Scan(&productID, &createdAt, &name, &category, &price)
			if err != nil {
				http.Error(w, "Failed to scan row", http.StatusInternalServerError)
				return
			}

			writer.Write([]string{
				strconv.Itoa(productID),
				name,
				category,
				strconv.FormatFloat(float64(price)/100, 'f', 2, 64),
				createdAt.Format("2006-01-02"),
			})
		}

		writer.Flush()
		csvFile.Seek(0, io.SeekStart)

		// Создаем ZIP-архив
		zipFile, err := os.CreateTemp("", "data-*.zip")
		if err != nil {
			http.Error(w, "Failed to create ZIP file", http.StatusInternalServerError)
			return
		}
		defer os.Remove(zipFile.Name())

		zipWriter := zip.NewWriter(zipFile)
		defer zipWriter.Close()

		csvInZip, err := zipWriter.Create("data.csv")
		if err != nil {
			http.Error(w, "Failed to create file in ZIP archive", http.StatusInternalServerError)
			return
		}

		_, err = io.Copy(csvInZip, csvFile)
		if err != nil {
			http.Error(w, "Failed to copy CSV to ZIP", http.StatusInternalServerError)
			return
		}

		// Отправляем архив клиенту
		w.Header().Set("Content-Type", "application/zip")
		w.Header().Set("Content-Disposition", "attachment; filename=data.zip")
		http.ServeFile(w, r, zipFile.Name())
	}
}
