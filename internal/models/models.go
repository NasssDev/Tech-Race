package models

import "time"

type Session struct {
	ID          int       `db:"id"`
	StartDate   time.Time `db:"start_time"`
	EndDate     time.Time `db:"end_time"`
	isAutopilot bool      `db:"is_autopilot"`
}

type Collision struct {
	ID          int       `db:"id"`
	Distance    float64   `db:"distance"`
	IsCollision bool      `db:"is_collision"`
	Timestamp   time.Time `db:"timestamp"`
	idSession   int       `db:"id_session"`
}

type Speed struct {
	ID        int       `db:"id"`
	Speed     string    `db:"speed"`
	Timestamp time.Time `db:"timestamp"`
	idSession int       `db:"id_session"`
}

type LineTracking struct {
	ID                int `db:"id"`
	LineTrackingValue int `db:"line_tracking_value"`
	idSession         int `db:"id_session"`
}

type Video struct {
	ID        int    `db:"id"`
	VideoURL  string `db:"video_url"`
	idSession int    `db:"id_session"`
}

type DatabaseInterface interface {
	GetAll() ([]Session, error)
	StartSession(isAutopilot bool) error
	StopSession() error
	IsSessionActive() (bool, error)
	InsertTrackData(data LineTracking) error
	InsertSonarData(data Collision) error
}
