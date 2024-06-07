package database

import "hetic/tech-race/internal/models"

func (d *Database) GetAllVideos() ([]models.Video, error) {
	// Implement logic to fetch all sessions from the database
	videos := []models.Video{}
	query := "SELECT * FROM videos"
	err := d.db.Select(&videos, query)
	if err != nil {
		return nil, err
	}
	return videos, nil
}
