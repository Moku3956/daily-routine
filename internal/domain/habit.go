package domain

type Habit struct {
	HabitId   int
	UserId    string
	HabitName string
	Category  string
	Value     int
	Unit      string
	PCost     int
	MCost     int
	Must      bool
}

type HabitRepository interface {
	Save(habit *Habit) error
	Update(habit *Habit) error
	Delete(id int) error
}
