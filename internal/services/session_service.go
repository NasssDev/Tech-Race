package services

import (
	"fmt"
	"hetic/tech-race/internal/models"
	"hetic/tech-race/internal/mqtt"
	"time"
)

type SessionService struct {
	db   models.DatabaseInterface
	info models.SessionInfo
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
	mqttClient := mqtt.NewMQTTClient(s.db)
	_ = mqttClient.ConnectAndSubscribe()
	fmt.Println("Session started")
	return nil
}

func (s *SessionService) Stop() error {
	err := s.db.StopSession()
	if err != nil {
		return err
	}
	mqttClient := mqtt.NewMQTTClient(s.db)
	mqttClient.Disconnect()
	fmt.Println("Session stopped")
	return nil
}
func (s *SessionService) IsSessionActive() (bool, error) {
	return s.db.IsSessionActive()
}

func (s *SessionService) GetAllSessionInfo() ([]models.SessionInfo, error) {
	sessions, err := s.db.GetAll()
	if err != nil {
		return nil, err
	}

	var sessionInfos []models.SessionInfo
	for _, session := range sessions {
		collisions, err := s.db.GetCollisionsBySessionID(session.ID)
		if err != nil {
			return nil, err
		}

		tracks, err := s.db.GetTracksBySessionID(session.ID)
		if err != nil {
			return nil, err
		}

		collisionInfo := models.CollisionInfo{
			Count:      len(collisions),
			Timestamps: make([]time.Time, len(collisions)),
		}
		for i, collision := range collisions {
			collisionInfo.Timestamps[i] = collision.Timestamp
		}

		trackInfo := models.TrackInfo{
			Count:      len(tracks),
			Timestamps: make([]time.Time, len(tracks)),
		}
		for i, collision := range collisions {
			collisionInfo.Timestamps[i] = collision.Timestamp // Make sure collision is of type Collision
		}

		// Format the start and end dates
		startDate := session.StartDate.Format("02_01_2006")
		endDate := session.EndDate.Format("02_01_2006")

		// Calculate the duration of the race
		duration := session.EndDate.Sub(session.StartDate)
		hours := int(duration.Hours())
		minutes := int(duration.Minutes()) % 60
		seconds := int(duration.Seconds()) % 60
		durationStr := fmt.Sprintf("%d_%d_%d", hours, minutes, seconds)

		sessionInfo := models.SessionInfo{
			ID:          session.ID,
			StartDate:   startDate,
			EndDate:     endDate,
			Duration:    durationStr,
			IsAutopilot: session.IsAutopilot,
			Collisions:  []models.CollisionInfo{collisionInfo},
			Tracks:      []models.TrackInfo{trackInfo},
		}

		sessionInfos = append(sessionInfos, sessionInfo)
	}

	return sessionInfos, nil
}
