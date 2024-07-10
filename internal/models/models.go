package models

import "time"

// Session struct for database
type Session struct {
	ID          int       `db:"id"`
	StartDate   time.Time `db:"start_time"`
	EndDate     time.Time `db:"end_time"`
	IsAutopilot bool      `db:"is_autopilot"`
}

type Collision struct {
	ID          int       `db:"id"`
	Distance    float64   `db:"distance"`
	IsCollision bool      `db:"is_collision"`
	Timestamp   time.Time `db:"timestamp"`
	IDSession   int       `db:"id_session"`
}

type Speed struct {
	ID        int       `db:"id"`
	Speed     string    `db:"speed"`
	Timestamp time.Time `db:"timestamp"`
	IDSession int       `db:"id_session"`
}

type LineTracking struct {
	ID                int       `db:"id"`
	LineTrackingValue int       `db:"line_tracking_value"`
	IDSession         int       `db:"id_session"`
	Timestamp         time.Time `db:"timestamp"`
}

type Video struct {
	ID        int    `db:"id"`
	VideoURL  string `db:"video_url"`
	IDSession int    `db:"id_session"`
}

// SessionInfo struct for json api
type SessionInfo struct {
	ID          int             `json:"id"`
	StartDate   string          `json:"start_time"`
	EndDate     string          `json:"end_time"`
	Duration    string          `json:"duration"`
	IsAutopilot bool            `json:"is_autopilot"`
	Videos      VideoInfo       `json:"videos"`
	Collisions  []CollisionInfo `json:"collisions"`
	Tracks      []TrackInfo     `json:"tracks"`
}

type CollisionInfo struct {
	Count      int       `json:"count"`
	Distances  []float64 `json:"distances"`
	Timestamps []string  `json:"timestamps"`
}

type TrackInfo struct {
	Count              int      `json:"count"`
	LineTrackingValues []int    `json:"line_tracking_values"`
	Timestamps         []string `json:"timestamps"`
}
type VideoInfo struct {
	VideoURLs []string `json:"video_urls"`
}

type DatabaseInterface interface {
	GetAll() ([]Session, error)
	StartSession(timeStamp time.Time, isAutopilot bool) error
	StopSession(timeStamp time.Time) error
	IsSessionActive() (bool, error)
	InsertTrackData(data LineTracking) error
	InsertSonarData(data Collision) error
	GetCurrentSessionID() (int, error)
	GetCollisionsBySessionID(sessionID int) ([]Collision, error)
	GetTracksBySessionID(sessionID int) ([]LineTracking, error)
	GetVideosBySessionID(sessionID int) ([]Video, error)
}
