package services

import (
	"fmt"
	"hetic/tech-race/internal/models"
	"hetic/tech-race/internal/mqtt"
	"time"
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
	timeStamp := time.Now()
	err := s.db.StartSession(timeStamp, isAutopilot)
	if err != nil {
		return err
	}
	mqttClient := mqtt.NewMQTTClient(s.db)
	_ = mqttClient.ConnectAndSubscribe()
	fmt.Println("Session started")
	return nil
}

func (s *SessionService) Stop() error {
	timeStamp := time.Now()
	err := s.db.StopSession(timeStamp)
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
		startDate := session.StartDate.Format("02/01/2006 - 15:04:05")
		endDate := session.EndDate.Format("02/01/2006 - 15:04:05")
		// Calculate the duration of the race
		duration := session.EndDate.Sub(session.StartDate)
		hours := int(duration.Hours())
		minutes := int(duration.Minutes()) % 60
		seconds := int(duration.Seconds()) % 60
		durationStr := fmt.Sprintf("%d:%d:%d", hours, minutes, seconds)
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
			Timestamps: make([]string, len(collisions)),
		}
		for i, collision := range collisions {
			timingCollision := collision.Timestamp.Sub(session.StartDate)
			hours := int(timingCollision.Hours())
			minutes := int(timingCollision.Minutes()) % 60
			seconds := int(timingCollision.Seconds()) % 60
			collisionInfo.Timestamps[i] = fmt.Sprintf("%d:%d:%d", hours, minutes, seconds)

		}

		trackInfo := models.TrackInfo{
			Count:      len(tracks),
			Timestamps: make([]string, len(tracks)),
		}
		for i, track := range tracks {
			timingTrack := track.Timestamp.Sub(session.StartDate)
			hours := int(timingTrack.Hours())
			minutes := int(timingTrack.Minutes()) % 60
			seconds := int(timingTrack.Seconds()) % 60
			trackInfo.Timestamps[i] = fmt.Sprintf("%d:%d:%d", hours, minutes, seconds)
		}

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
