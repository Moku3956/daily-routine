package domain

type Habit struct {
	HabitId     int
	UserId      int
	HabitName   string
	Category    string
	TargetValue int
	Unit        string
	PCost       int
	MCost       int
	Must        bool
}

type HabitRepository interface {
	Save(habit *Habit) error
	Update(habit *Habit) error
	Delete(id int) error
	FindAll(userId int) ([]*Habit, error)
	FindById(id int) (*Habit, error)
}
