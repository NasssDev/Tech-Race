package database

import "hetic/tech-race/internal/models"

func (d *Database) GetAll() ([]models.Session, error) {
	sessions := []models.Session{}
	query := "SELECT * FROM Session"
	err := d.db.Select(&sessions, query)
	if err != nil {
		return nil, err
	}
	return sessions, nil
}
func (d *Database) StartSession(isAutopilot bool) error {
	query := `INSERT INTO Session (start_time, end_time, is_autopilot) VALUES (NOW(), NULL, $1)`
	_, err := d.db.Exec(query, isAutopilot)
	if err != nil {
		return err
	}
	return nil
}

func (d *Database) StopSession() error {
	query := `UPDATE Session SET end_time = NOW() WHERE end_time IS NULL`
	_, err := d.db.Exec(query)
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
	query := `INSERT INTO LineTracking (line_tracking_value) VALUES ($1)`
	_, err := d.db.Exec(query, data.LineTrackingValue)
	if err != nil {
		return err
	}
	return nil
}

func (d *Database) InsertSonarData(data models.Collision) error {
	query := `INSERT INTO Collision (distance, is_collision, timestamp) VALUES ($1, $2, $3)`
	_, err := d.db.Exec(query, data.Distance, data.IsCollision, data.Timestamp)
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
