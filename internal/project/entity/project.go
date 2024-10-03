package entity

import (
	"github.com/google/uuid"
	"time"
)

type ProjectId uuid.UUID

type Project struct {
	Id              ProjectId
	Name            string
	Description     string
	PreviewImageUrl string
	LinkedLink      string
	GithubLink      string
	Date            time.Time
	Client          string
	Categories      []*ProjectCategory
	TechStack       []*ProjectTech
	ImagesLinks     []string
}

func (u ProjectId) String() string {
	return uuid.UUID(u).String()
}

func (u *ProjectId) Scan(src interface{}) error {
	us := (*uuid.UUID)(u)
	return us.Scan(src)
}
