package usecase

import (
	"errors"
	"fmt"
	"time"

	"github.com/Moku3956/daily-routine/internal/domain"
)

var (
	ErrConditionAlreadyExists = errors.New("本日のコンディションはすでに入力されています")
	ErrConditionNotFound      = errors.New("本日のコンディションが見つかりません")
)

type ConditionUsecase struct {
	repo domain.ConditionRepository
}

func NewConditionUsecase(repo domain.ConditionRepository) *ConditionUsecase {
	return &ConditionUsecase{repo: repo}
}

func (u *ConditionUsecase) InputCondition(userId, pCondition, mCondition int) (*domain.Condition, error) {
	today := today()

	existing, err := u.repo.FindByUserIdAndDate(userId, today)
	if err != nil {
		return nil, fmt.Errorf("コンディション検索に失敗: %w", err)
	}
	if existing != nil {
		return nil, ErrConditionAlreadyExists
	}

	c := &domain.Condition{
		UserId:     userId,
		PCondition: pCondition,
		MCondition: mCondition,
		Rank:       domain.CalcConditionRank(pCondition, mCondition),
		Date:       today,
	}
	if err := u.repo.Save(c); err != nil {
		return nil, fmt.Errorf("コンディション登録に失敗: %w", err)
	}
	return c, nil
}

func (u *ConditionUsecase) GetTodayCondition(userId int) (*domain.Condition, error) {
	c, err := u.repo.FindByUserIdAndDate(userId, today())
	if err != nil {
		return nil, fmt.Errorf("コンディション取得に失敗: %w", err)
	}
	if c == nil {
		return nil, ErrConditionNotFound
	}
	return c, nil
}

func today() time.Time {
	t := time.Now()
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}
