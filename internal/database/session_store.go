package database

import (
	"hetic/tech-race/internal/models"
	"time"
)

func (d *Database) GetAll() ([]models.Session, error) {
	sessions := []models.Session{}
	query := "SELECT * FROM Session"
	err := d.db.Select(&sessions, query)
	if err != nil {
		return nil, err
	}
	return sessions, nil
}
func (d *Database) StartSession(timeStamp time.Time, isAutopilot bool) error {
	query := `INSERT INTO Session (start_time, end_time, is_autopilot) VALUES ($1, NULL, $2)`
	_, err := d.db.Exec(query, timeStamp, isAutopilot)
	if err != nil {
		return err
	}
	return nil
}

func (d *Database) StopSession(timeStamp time.Time) error {
	query := `UPDATE Session SET end_time = $1 WHERE end_time IS NULL`
	_, err := d.db.Exec(query, timeStamp)
	if err != nil {
		return err
	}
	return nil
}
func (d *Database) IsSessionActive() (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM Session WHERE end_time IS NULL`
	err := d.db.Get(&count, query)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
func (d *Database) InsertTrackData(data models.LineTracking) error {
	query := `INSERT INTO LineTracking (line_tracking_value , id_session, timestamp) VALUES ($1, $2, $3)`
	_, err := d.db.Exec(query, data.LineTrackingValue, data.IDSession, data.Timestamp)
	if err != nil {
		return err
	}
	return nil
}

func (d *Database) InsertSonarData(data models.Collision) error {
	query := `INSERT INTO Collision (distance, is_collision, timestamp, id_session) VALUES ($1, $2, $3, $4)`
	_, err := d.db.Exec(query, data.Distance, data.IsCollision, data.Timestamp, data.IDSession)
	if err != nil {
		return err
	}
	return nil
}
func (d *Database) GetCurrentSessionID() (int, error) {
	var id int
	query := `SELECT id FROM Session WHERE end_time IS NULL`
	err := d.db.Get(&id, query)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (d *Database) GetLastSessionID() (int, error) {
	var id int
	query := `SELECT id FROM Session ORDER BY id DESC LIMIT 1`
	err := d.db.Get(&id, query)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (d *Database) GetCollisionsBySessionID(sessionID int) ([]models.Collision, error) {
	collisions := []models.Collision{}
	query := "SELECT * FROM Collision WHERE id_session = $1"
	err := d.db.Select(&collisions, query, sessionID)
	if err != nil {
		return nil, err
	}
	return collisions, nil
}

func (d *Database) GetTracksBySessionID(sessionID int) ([]models.LineTracking, error) {
	tracks := []models.LineTracking{}
	query := "SELECT * FROM LineTracking WHERE id_session = $1"
	err := d.db.Select(&tracks, query, sessionID)
	if err != nil {
		return nil, err
	}
	return tracks, nil
}

func (d *Database) InsertVideoData(data models.Video) error {
	query := `INSERT INTO video (video_url, id_session) VALUES ($1, $2)`
	_, err := d.db.Exec(query, data.VideoURL, data.IDSession)
	if err != nil {
		return err
	}
	return nil
}
