package domain

type User struct {
	UserId   int
	UserName string
	Password string
}

type Session struct {
	SessionId int
	UserId    int
	Token     string
}

type UserRepository interface {
	Save(user *User) error
	FindByUserName(userName string) (*User, error)
}

type SessionRepository interface {
	Save(session *Session) error
	FindByToken(token string) (*Session, error)
	DeleteByToken(token string) error
}
