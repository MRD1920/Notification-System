package DB

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool
var once sync.Once

func connectDB() *pgxpool.Pool {
	//db connection
	// dbURL := os.Getenv("DB_URL")
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	if dbURL == "" {
		log.Fatal("DB_URL is not set in .env")
	}

	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		log.Fatalf("Unable to parse DB URL: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
	}

	log.Println("Connected to DB!!!")
	return pool
}

func InitDBPool() {
	once.Do(func() {
		Pool = connectDB()
	})
}

func CloseDBPool() {
	if Pool != nil {
		Pool.Close()
	}
}
