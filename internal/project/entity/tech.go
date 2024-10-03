package entity

import "github.com/google/uuid"

type ProjectTechId uuid.UUID

type ProjectTech struct {
	Id   ProjectTechId
	Name string
}

func (u ProjectTechId) String() string {
	return uuid.UUID(u).String()
}

func (u *ProjectTechId) Scan(src interface{}) error {
	us := (*uuid.UUID)(u)
	return us.Scan(src)
}
