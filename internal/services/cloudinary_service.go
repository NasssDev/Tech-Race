package services

import (
	"fmt"
	"hetic/tech-race/internal/models"
)

type CloudinaryService struct {
	db models.Database
}

func NewCloudinaryService(db models.Database) *CloudinaryService {
	return &CloudinaryService{db: db}
}

func (s *CloudinaryService) GetAllVideos() {
	sessions, err := s.db.GetAllVideos()
	if err != nil {
		return
	}
	fmt.Println(sessions)
}