package middleware

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Moku3956/daily-routine/internal/domain"
)

type contextKey string

const UserIDKey contextKey = "user_id"

func GetUserID(ctx context.Context) (int, bool) {
	id, ok := ctx.Value(UserIDKey).(int)
	return id, ok
}

type sessionFinder interface {
	FindByToken(token string) (*domain.Session, error)
}

func Auth(sessionRepo sessionFinder) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("session_token")
			if err != nil || cookie.Value == "" {
				writeUnauthorized(w, "認証が必要です")
				return
			}
			session, err := sessionRepo.FindByToken(cookie.Value)
			if err != nil || session == nil {
				writeUnauthorized(w, "セッションが無効です")
				return
			}
			ctx := context.WithValue(r.Context(), UserIDKey, session.UserId)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func writeUnauthorized(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(map[string]any{
		"status":  401,
		"error":   "Unauthorized",
		"message": message,
	})
}
