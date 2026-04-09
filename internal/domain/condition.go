package domain

import "time"

type ConditionRank string

const (
	RankGood   ConditionRank = "good"
	RankNormal ConditionRank = "normal"
	RankBad    ConditionRank = "bad"
)

type Condition struct {
	ConditionId int
	UserId      int
	PCondition  int
	MCondition  int
	Rank        ConditionRank
	Date        time.Time
}

// CalcConditionRank は合計値でランクを決定する（8以上:good, 5-7:normal, 4以下:bad）
func CalcConditionRank(p, m int) ConditionRank {
	sum := p + m
	switch {
	case sum >= 8:
		return RankGood
	case sum >= 5:
		return RankNormal
	default:
		return RankBad
	}
}

type ConditionRepository interface {
	Save(condition *Condition) error
	FindByUserIdAndDate(userId int, date time.Time) (*Condition, error)
}
