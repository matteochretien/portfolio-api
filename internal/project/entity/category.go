package entity

import "github.com/google/uuid"

type ProjectCategoryId uuid.UUID

type ProjectCategory struct {
	Id   ProjectCategoryId
	Name string
}

func (u ProjectCategoryId) String() string {
	return uuid.UUID(u).String()
}

func (u *ProjectCategoryId) Scan(src interface{}) error {
	us := (*uuid.UUID)(u)
	return us.Scan(src)
}
