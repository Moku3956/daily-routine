package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Moku3956/daily-routine/internal/adapter/handler"
	"github.com/Moku3956/daily-routine/internal/adapter/repository"
	"github.com/Moku3956/daily-routine/internal/usecase"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func run() error {
	// .envの読み込み
	if err := godotenv.Load(); err != nil {
		// Docker環境では.envがないのは正常
		// ローカル実行時のみ警告を出す
		if os.Getenv("DATABASE_URL") == "" {
			log.Printf("warning: .envファイルを読み込めませんでした: %v", err)
		}
	}

	// DB接続
	conn := os.Getenv("DATABASE_URL")
	if conn == "" {
		return fmt.Errorf("DATABASE_URL環境変数が設定されていません")
	}
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return fmt.Errorf("DB接続の初期化に失敗しました: %w", err)
	}
	defer db.Close()

	// setup.sqlの読み込み
	schemaBytes, err := os.ReadFile("internal/adapter/repository/setup.sql")
	if err != nil {
		return fmt.Errorf("setup.sql の読み込みに失敗しました: %w", err)
	}
	schemaSQL := strings.TrimSpace(string(schemaBytes))
	if schemaSQL == "" {
		return fmt.Errorf("setup.sql が空です")
	}

	// Contextを使って、デバックしやすくする
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// DBとの通信確認
	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("DBへの接続に失敗しました: %w", err)
	}
	fmt.Println("DB接続成功")

	// 読み込んだsetup.sqlを使って、テーブル作成
	if _, err := db.ExecContext(ctx, schemaSQL); err != nil {
		return fmt.Errorf("habits テーブル作成に失敗しました: %w", err)
	}

	habitRepo := repository.NewSqlHabitRepository(db)
	habitUsecase := usecase.NewHabitUsecase(habitRepo)
	habitHandler := handler.NewHabitHttpHandler(habitUsecase)

	http.HandleFunc("/habit", habitHandler.Create)
	fmt.Println("サーバー起動中 :8080")
	return http.ListenAndServe(":8080", nil)
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("%v", err)
	}
}
