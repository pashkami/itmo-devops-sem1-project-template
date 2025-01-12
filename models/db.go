package models

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)


func ConnectDB(dataSourceName string) (*sql.DB, error) {
    db, err := sql.Open("postgres", dataSourceName)
    if err != nil {
        return nil, err
    }

    if err := db.Ping(); err != nil {
        return nil, err
    }

    // Создание таблицы
    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS prices (
            id SERIAL PRIMARY KEY,
            created_at DATE NOT NULL,
            name VARCHAR(255) NOT NULL,
            category VARCHAR(255) NOT NULL,
            price DECIMAL(10, 2) NOT NULL
        )
    `)
	fmt.Println("Connected to the database!")
    return db, err
}
