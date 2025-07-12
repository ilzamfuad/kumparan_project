package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func BuildDB() *sql.DB {

	DbHost := os.Getenv("DB_HOST")
	DbUser := os.Getenv("DB_USER")
	DbPassword := os.Getenv("DB_PASSWORD")
	DbName := os.Getenv("DB_NAME")
	DbPort := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", DbUser, DbPassword, DbHost, DbPort, DbName)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("Failed to open DB: %v", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(30 * time.Minute)

	fmt.Println("Database connection established")
	return db
}

// If you are using GORM, you can uncomment the following lines and use GORM instead of sql.DB
// var DB *gorm.DB

// func BuildDB() *gorm.DB {
// 	DbHost := os.Getenv("DB_HOST")
// 	DbUser := os.Getenv("DB_USER")
// 	DbPassword := os.Getenv("DB_PASSWORD")
// 	DbName := os.Getenv("DB_NAME")
// 	DbPort := os.Getenv("DB_PORT")
// 	sqlCfg := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
// 		DbHost,
// 		DbPort,
// 		DbUser,
// 		DbPassword,
// 		DbName,
// 	)

// 	db, err := gorm.Open(postgres.Open(sqlCfg), &gorm.Config{})

// 	if err != nil {
// 		fmt.Println("Cannot connect to database postgres")
// 		log.Fatal("connection error:", err)
// 	}

// 	sqlDB, err := db.DB()
// 	if err != nil {
// 		panic(err)
// 	}

// 	sqlDB.SetMaxIdleConns(5)
// 	sqlDB.SetMaxOpenConns(10)
// 	sqlDB.SetConnMaxLifetime(15 * time.Minute)

// 	return db
// }
