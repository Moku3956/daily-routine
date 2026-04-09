package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/Moku3956/daily-routine/internal/domain"
	"github.com/Moku3956/daily-routine/internal/usecase"
)

type userUsecase interface {
	Register(userName, password string) (*domain.User, error)
	Login(userName, password string) (*domain.Session, error)
	Logout(token string) error
}

type AuthHttpHandler struct {
	usecase userUsecase
}

func NewAuthHttpHandler(u userUsecase) *AuthHttpHandler {
	return &AuthHttpHandler{usecase: u}
}

type registerRequest struct {
	UserName        string `json:"user_name"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
}

type loginRequest struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

func (h *AuthHttpHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "リクエストの形式が不正です")
		return
	}
	if req.UserName == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "ユーザー名とパスワードは必須です")
		return
	}
	if req.Password != req.PasswordConfirm {
		writeError(w, http.StatusUnprocessableEntity, "パスワードが一致しません")
		return
	}

	user, err := h.usecase.Register(req.UserName, req.Password)
	if err != nil {
		if errors.Is(err, usecase.ErrUserAlreadyExists) {
			writeError(w, http.StatusConflict, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "ユーザー登録に失敗しました")
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"user_id":   user.UserId,
		"user_name": user.UserName,
	})
}

func (h *AuthHttpHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "リクエストの形式が不正です")
		return
	}
	if req.UserName == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "ユーザー名とパスワードは必須です")
		return
	}

	session, err := h.usecase.Login(req.UserName, req.Password)
	if err != nil {
		if errors.Is(err, usecase.ErrUserNotFound) {
			writeError(w, http.StatusNotFound, "ユーザーが見つかりません")
			return
		}
		if errors.Is(err, usecase.ErrInvalidPassword) {
			writeError(w, http.StatusUnauthorized, "ログインに失敗しました")
			return
		}
		writeError(w, http.StatusInternalServerError, "ログインに失敗しました")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    session.Token,
		HttpOnly: true,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
	})
	writeJSON(w, http.StatusCreated, map[string]any{
		"session_id": session.SessionId,
	})
}

func (h *AuthHttpHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		writeError(w, http.StatusUnauthorized, "セッションが見つかりません")
		return
	}

	if err := h.usecase.Logout(cookie.Value); err != nil {
		writeError(w, http.StatusInternalServerError, "ログアウトに失敗しました")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:   "session_token",
		Value:  "",
		MaxAge: -1,
		Path:   "/",
	})
	w.WriteHeader(http.StatusNoContent)
}
