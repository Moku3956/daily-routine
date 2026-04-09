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
	"github.com/Moku3956/daily-routine/internal/adapter/middleware"
	"github.com/Moku3956/daily-routine/internal/adapter/repository"
	"github.com/Moku3956/daily-routine/internal/usecase"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func run() error {
	if err := godotenv.Load(); err != nil {
		if os.Getenv("DATABASE_URL") == "" {
			log.Printf("warning: .envファイルを読み込めませんでした: %v", err)
		}
	}

	conn := os.Getenv("DATABASE_URL")
	if conn == "" {
		return fmt.Errorf("DATABASE_URL環境変数が設定されていません")
	}
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return fmt.Errorf("DB接続の初期化に失敗しました: %w", err)
	}
	defer db.Close()

	schemaBytes, err := os.ReadFile("internal/adapter/repository/setup.sql")
	if err != nil {
		return fmt.Errorf("setup.sql の読み込みに失敗しました: %w", err)
	}
	schemaSQL := strings.TrimSpace(string(schemaBytes))
	if schemaSQL == "" {
		return fmt.Errorf("setup.sql が空です")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("DBへの接続に失敗しました: %w", err)
	}
	fmt.Println("DB接続成功")

	if _, err := db.ExecContext(ctx, schemaSQL); err != nil {
		return fmt.Errorf("テーブル作成に失敗しました: %w", err)
	}

	// リポジトリ
	userRepo := repository.NewSqlUserRepository(db)
	sessionRepo := repository.NewSqlSessionRepository(db)
	habitRepo := repository.NewSqlHabitRepository(db)
	conditionRepo := repository.NewSqlConditionRepository(db)
	logRepo := repository.NewSqlLogRepository(db)

	// ユースケース
	userUC := usecase.NewUserUsecase(userRepo, sessionRepo)
	habitUC := usecase.NewHabitUsecase(habitRepo, conditionRepo)
	conditionUC := usecase.NewConditionUsecase(conditionRepo)
	logUC := usecase.NewLogUsecase(logRepo, habitRepo)
	statsUC := usecase.NewStatsUsecase(logRepo)

	// ハンドラー
	authHandler := handler.NewAuthHttpHandler(userUC)
	habitHandler := handler.NewHabitHttpHandler(habitUC)
	conditionHandler := handler.NewConditionHttpHandler(conditionUC)
	logHandler := handler.NewLogHttpHandler(logUC)
	statsHandler := handler.NewStatsHttpHandler(statsUC)

	// 認証ミドルウェア
	auth := middleware.Auth(sessionRepo)
	protect := func(h http.HandlerFunc) http.Handler { return auth(h) }

	// ルーティング
	mux := http.NewServeMux()

	// 認証不要
	mux.HandleFunc("POST /register", authHandler.Register)
	mux.HandleFunc("POST /session", authHandler.Login)
	mux.HandleFunc("DELETE /session", authHandler.Logout)

	// 認証必要
	mux.Handle("POST /condition", protect(conditionHandler.Create))
	mux.Handle("GET /condition", protect(conditionHandler.Get))

	mux.Handle("POST /habits", protect(habitHandler.Create))
	mux.Handle("GET /habits", protect(habitHandler.List))
	mux.Handle("GET /habits/ordered", protect(habitHandler.GetOrdered))
	mux.Handle("PATCH /habits/{id}", protect(habitHandler.Update))
	mux.Handle("DELETE /habits/{id}", protect(habitHandler.Delete))

	mux.Handle("POST /habits/logs", protect(logHandler.Create))
	mux.Handle("GET /habits/logs", protect(logHandler.List))

	mux.Handle("GET /stats", protect(statsHandler.Get))

	fmt.Println("サーバー起動中 :8080")
	return http.ListenAndServe(":8080", mux)
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("%v", err)
	}
}
