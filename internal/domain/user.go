package domain

type User struct {
	ID       int64
	Email    string
	Phone    string
	Password string

	Ctime int64
}
