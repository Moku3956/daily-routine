package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Moku3956/daily-routine/internal/adapter/repository"
	"github.com/Moku3956/daily-routine/internal/domain"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func run() error {
	if err := godotenv.Load(); err != nil {
		log.Printf("warning: .envファイルを読み込めませんでした: %v", err)
	}

	password := os.Getenv("POSTGRES_PASSWORD")

	conn := fmt.Sprintf("user=postgres password=%s dbname=postgres host=127.0.0.1 port=5432 sslmode=disable", password)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return fmt.Errorf("DB接続の初期化に失敗しました: %w", err)
	}

	schemaBytes, err := os.ReadFile("internal/adapter/repository/setup.sql")
	if err != nil {
		return fmt.Errorf("setup.sql の読み込みに失敗しました: %w", err)
	}
	schemaSQL := strings.TrimSpace(string(schemaBytes))
	if schemaSQL == "" {
		return fmt.Errorf("setup.sql が空です")
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("DBへの接続に失敗しました: %w", err)
	}
	fmt.Println("DB接続成功")

	if _, err := db.ExecContext(ctx, schemaSQL); err != nil {
		return fmt.Errorf("habits テーブル作成に失敗しました: %w", err)
	}

	habitRepo := repository.NewSqlHabitRepository(db)
	newHabit := domain.Habit{
		UserId:    1,
		HabitName: "朝の散歩",
		Category:  "健康",
		Value:     30,
		Unit:      "分",
		PCost:     2,
		MCost:     1,
		Must:      true,
	}
	if err := habitRepo.Save(&newHabit); err != nil {
		return fmt.Errorf("習慣登録に失敗しました: %w", err)
	}
	fmt.Println("習慣登録成功")

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("習慣の登録に失敗しました: %v", err)
	}
}
