package models

import "time"

type Session struct {
	StartDate   time.Time
	EndDate     time.Time
	isAutopilot bool
}

type Collision struct {
	Distance    string
	IsCollision bool
	Timestamp   time.Time
}

type Speed struct {
	Speed     string
	Timestamp time.Time
}

type LineTracking struct {
	LineTrackingValue int
}

type Video struct {
	VideoURL string
}

type DatabaseInterface interface {
	GetAll() ([]Session, error)
}
