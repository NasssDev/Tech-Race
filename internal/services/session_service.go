package services

import (
	"fmt"
	"hetic/tech-race/internal/models"
)

type SessionService struct {
	db models.DatabaseInterface
}

func NewSessionService(db models.DatabaseInterface) *SessionService {
	return &SessionService{db: db}
}

func (s *SessionService) GetAll() ([]models.Session, error) {
	sessions, err := s.db.GetAll()
	if err != nil {
		return nil, err
	}
	fmt.Println(sessions)
	return sessions, nil
}
func (s *SessionService) Start(isAutopilot bool) error {
	err := s.db.StartSession(isAutopilot)
	if err != nil {
		return err
	}
	fmt.Println("Session started")
	return nil
}

func (s *SessionService) Stop() error {
	err := s.db.StopSession()
	if err != nil {
		return err
	}
	fmt.Println("Session stopped")
	return nil
}
func (s *SessionService) IsSessionActive() (bool, error) {
	return s.db.IsSessionActive()
}
