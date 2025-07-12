package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	DbHost := os.Getenv("DB_HOST")
	DbUser := os.Getenv("DB_USER")
	DbPassword := os.Getenv("DB_PASSWORD")
	DbName := os.Getenv("DB_NAME")
	DbPort := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", DbUser, DbPassword, DbHost, DbPort, DbName)
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	authors := []struct {
		ID   string
		Name string
	}{
		{ID: uuid.New().String(), Name: "Author One"},
		{ID: uuid.New().String(), Name: "Author Two"},
		{ID: uuid.New().String(), Name: "Author Three"},
	}

	for _, author := range authors {
		_, err := db.Exec("INSERT INTO authors (id, name) VALUES ($1, $2)", author.ID, author.Name)
		if err != nil {
			log.Fatalf("Failed to insert author %s: %v", author.Name, err)
		}
		fmt.Printf("Inserted author: %s\n", author.Name)
	}

	fmt.Println("Seeding completed successfully.")
}
