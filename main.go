package main

import (
	"project_sem/api"
	"project_sem/models"
	"net/http"
)

func main() {
	// Подключение к базе данных
	db := models.ConnectDB()
	defer db.Close()

	// Маршруты
	http.HandleFunc("/upload", api.UploadPricesHandler(db))
	http.HandleFunc("/download", api.DownloadPricesHandler(db))

	// Запуск сервера
	serverAddress := ":8080"
	println("Server is running on", serverAddress)
	if err := http.ListenAndServe(serverAddress, nil); err != nil {
		panic(err)
	}
}
