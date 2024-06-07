package database

import "hetic/tech-race/internal/models"

func (d *Database) GetAll() ([]models.Session, error) {
	// Implement logic to fetch all sessions from the database
	sessions := []models.Session{}
	query := "SELECT * FROM sessions"
	err := d.db.Select(&sessions, query)
	if err != nil {
		return nil, err
	}
	return sessions, nil
}
