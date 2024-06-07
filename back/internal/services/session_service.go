package services

import (
	"fmt"
	"hetic/tech-race/internal/models"
)

type SessionService struct {
	db models.Database
}

func NewSessionService(db models.Database) *SessionService {
	return &SessionService{db: db}
}

func (s *SessionService) GetAll() {
	sessions, err := s.db.GetAll()
	if err != nil {
		return
	}
	fmt.Println(sessions)
}
