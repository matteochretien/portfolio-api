package entity

import "github.com/google/uuid"

type UserId uuid.UUID

type User struct {
	Id          UserId
	FirstName   string
	LastName    string
	Email       string
	Password    string
	Permissions []string
}

func (u UserId) String() string {
	return uuid.UUID(u).String()
}

func (u *UserId) Scan(src interface{}) error {
	us := (*uuid.UUID)(u)
	return us.Scan(src)
}
