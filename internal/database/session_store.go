package database

import "hetic/tech-race/internal/models"

func (d *Database) GetAll() ([]models.Session, error) {
	sessions := []models.Session{}
	query := "SELECT * FROM sessions"
	err := d.db.Select(&sessions, query)
	if err != nil {
		return nil, err
	}
	return sessions, nil
}
