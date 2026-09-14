package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf(".env読み込み失敗: %v", err)
	}

	db, err := sql.Open("pgx", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("DB接続失敗: %v", err)
	}
	defer db.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("dialect設定失敗: %v", err)
	}

	if err := goose.Up(db, "migration"); err != nil {
		log.Fatalf("マイグレーション適用失敗: %v", err)
	}

	log.Println("マイグレーション完了")
}
