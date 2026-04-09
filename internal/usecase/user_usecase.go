package usecase

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/Moku3956/daily-routine/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserAlreadyExists = errors.New("ユーザー名はすでに使用されています")
	ErrUserNotFound      = errors.New("ユーザーが見つかりません")
	ErrInvalidPassword   = errors.New("パスワードが正しくありません")
)

type UserUsecase struct {
	userRepo    domain.UserRepository
	sessionRepo domain.SessionRepository
}

func NewUserUsecase(userRepo domain.UserRepository, sessionRepo domain.SessionRepository) *UserUsecase {
	return &UserUsecase{userRepo: userRepo, sessionRepo: sessionRepo}
}

func (u *UserUsecase) Register(userName, password string) (*domain.User, error) {
	existing, err := u.userRepo.FindByUserName(userName)
	if err != nil {
		return nil, fmt.Errorf("ユーザー検索に失敗: %w", err)
	}
	if existing != nil {
		return nil, ErrUserAlreadyExists
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("パスワードハッシュ化に失敗: %w", err)
	}

	user := &domain.User{UserName: userName, Password: string(hashed)}
	if err := u.userRepo.Save(user); err != nil {
		return nil, fmt.Errorf("ユーザー登録に失敗: %w", err)
	}
	return user, nil
}

func (u *UserUsecase) Login(userName, password string) (*domain.Session, error) {
	user, err := u.userRepo.FindByUserName(userName)
	if err != nil {
		return nil, fmt.Errorf("ユーザー検索に失敗: %w", err)
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, ErrInvalidPassword
	}

	token, err := generateToken()
	if err != nil {
		return nil, fmt.Errorf("トークン生成に失敗: %w", err)
	}

	session := &domain.Session{UserId: user.UserId, Token: token}
	if err := u.sessionRepo.Save(session); err != nil {
		return nil, fmt.Errorf("セッション保存に失敗: %w", err)
	}
	return session, nil
}

func (u *UserUsecase) Logout(token string) error {
	if err := u.sessionRepo.DeleteByToken(token); err != nil {
		return fmt.Errorf("セッション削除に失敗: %w", err)
	}
	return nil
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
