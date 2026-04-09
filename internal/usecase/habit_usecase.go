package usecase

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"sort"
	"time"

	"github.com/Moku3956/daily-routine/internal/domain"
)

var (
	ErrHabitNotFound  = errors.New("習慣が見つかりません")
	ErrHabitForbidden = errors.New("この習慣を操作する権限がありません")
)

type HabitUsecase struct {
	habitRepo     domain.HabitRepository
	conditionRepo domain.ConditionRepository
}

func NewHabitUsecase(habitRepo domain.HabitRepository, conditionRepo domain.ConditionRepository) *HabitUsecase {
	return &HabitUsecase{habitRepo: habitRepo, conditionRepo: conditionRepo}
}

func (u *HabitUsecase) AddHabit(h *domain.Habit) error {
	if err := u.habitRepo.Save(h); err != nil {
		return fmt.Errorf("習慣登録に失敗: %w", err)
	}
	return nil
}

func (u *HabitUsecase) GetHabits(userId int) ([]*domain.Habit, error) {
	habits, err := u.habitRepo.FindAll(userId)
	if err != nil {
		return nil, fmt.Errorf("習慣一覧取得に失敗: %w", err)
	}
	return habits, nil
}

func (u *HabitUsecase) UpdateHabit(userId int, h *domain.Habit) (*domain.Habit, error) {
	existing, err := u.habitRepo.FindById(h.HabitId)
	if err != nil {
		return nil, fmt.Errorf("習慣取得に失敗: %w", err)
	}
	if existing == nil {
		return nil, ErrHabitNotFound
	}
	if existing.UserId != userId {
		return nil, ErrHabitForbidden
	}
	h.UserId = userId
	if err := u.habitRepo.Update(h); err != nil {
		return nil, fmt.Errorf("習慣更新に失敗: %w", err)
	}
	return h, nil
}

func (u *HabitUsecase) DeleteHabit(userId, habitId int) error {
	existing, err := u.habitRepo.FindById(habitId)
	if err != nil {
		return fmt.Errorf("習慣取得に失敗: %w", err)
	}
	if existing == nil {
		return ErrHabitNotFound
	}
	if existing.UserId != userId {
		return ErrHabitForbidden
	}
	if err := u.habitRepo.Delete(habitId); err != nil {
		return fmt.Errorf("習慣削除に失敗: %w", err)
	}
	return nil
}

// OrderedHabit はコンディションに応じて調整済みの習慣情報
type OrderedHabit struct {
	HabitId     int
	HabitName   string
	Category    string
	TargetValue int // 調整後の目標値
	Unit        string
	PCost       int
	MCost       int
	Must        bool
}

func (u *HabitUsecase) GetOrderedHabits(userId int) ([]*OrderedHabit, error) {
	condition, err := u.conditionRepo.FindByUserIdAndDate(userId, todayDate())
	if err != nil {
		return nil, fmt.Errorf("コンディション取得に失敗: %w", err)
	}
	if condition == nil {
		return nil, ErrConditionNotFound
	}

	habits, err := u.habitRepo.FindAll(userId)
	if err != nil {
		return nil, fmt.Errorf("習慣一覧取得に失敗: %w", err)
	}

	// 不調の場合はマスト習慣のみ
	var targets []*domain.Habit
	if condition.Rank == domain.RankBad {
		for _, h := range habits {
			if h.Must {
				targets = append(targets, h)
			}
		}
	} else {
		targets = habits
	}

	sort.SliceStable(targets, func(i, j int) bool {
		hi, hj := targets[i], targets[j]
		if hi.Must != hj.Must {
			return hi.Must // マスト優先
		}
		totalI := hi.PCost + hi.MCost
		totalJ := hj.PCost + hj.MCost
		if totalI != totalJ {
			if condition.Rank == domain.RankBad {
				return totalI < totalJ // 不調: コスト昇順
			}
			return totalI > totalJ // 好調/普通: コスト降順
		}
		// コストが同じ場合の tie-break
		return tieBreak(hi, hj, condition)
	})

	result := make([]*OrderedHabit, 0, len(targets))
	for _, h := range targets {
		displayValue := adjustValue(h.TargetValue, h.Must, condition.Rank)
		result = append(result, &OrderedHabit{
			HabitId:     h.HabitId,
			HabitName:   h.HabitName,
			Category:    h.Category,
			TargetValue: displayValue,
			Unit:        h.Unit,
			PCost:       h.PCost,
			MCost:       h.MCost,
			Must:        h.Must,
		})
	}
	return result, nil
}

// tieBreak はコストが同じ場合の順序を決める
func tieBreak(hi, hj *domain.Habit, c *domain.Condition) bool {
	if c.PCondition == c.MCondition {
		return rand.Intn(2) == 0
	}
	if c.PCondition > c.MCondition {
		if c.Rank == domain.RankBad {
			return hi.PCost < hj.PCost
		}
		return hi.PCost > hj.PCost
	}
	if c.Rank == domain.RankBad {
		return hi.MCost < hj.MCost
	}
	return hi.MCost > hj.MCost
}

// adjustValue はコンディションに応じて目標値を調整する
func adjustValue(base int, must bool, rank domain.ConditionRank) int {
	switch rank {
	case domain.RankGood:
		if must {
			return int(math.Round(float64(base) * 1.2))
		}
		return base
	case domain.RankBad:
		return int(math.Round(float64(base) * 0.5))
	default:
		return base
	}
}

func todayDate() time.Time {
	t := time.Now()
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}
