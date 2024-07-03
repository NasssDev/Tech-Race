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
	Distance    string    `db:"distance"`
	IsCollision bool      `db:"is_collision"`
	Timestamp   time.Time `db:"timestamp"`
}

type Speed struct {
	ID        int       `db:"id"`
	Speed     string    `db:"speed"`
	Timestamp time.Time `db:"timestamp"`
}

type LineTracking struct {
	ID                int `db:"id"`
	LineTrackingValue int `db:"line_tracking_value"`
}

type Video struct {
	ID       int    `db:"id"`
	VideoURL string `db:"video_url"`
}

type DatabaseInterface interface {
	GetAll() ([]Session, error)
	StartSession(isAutopilot bool) error
	StopSession() error
	IsSessionActive() (bool, error)
}
