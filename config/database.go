package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

var DB *sql.DB

func ConnectDB() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	// open connect (but dont ping)
	DB, err = sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal("Gagal membuka driver DB: ", err)
	}

	// test connection (ping)
	err = DB.Ping()
	if err != nil {
		log.Fatal("Gagal konek ke Database (Cek user/pass/port): ", err)
	}

	fmt.Println("🚀 Connected to PostgreSQL database successfully!")
}
