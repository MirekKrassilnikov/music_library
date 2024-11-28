package services

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/MirekKrassilnikov/music_library/dto"
	"github.com/MirekKrassilnikov/music_library/models"
)

type SongService struct {
	DB *sql.DB
}

// GetAllSongs обрабатывает бизнес-логику поиска песен
func (s *SongService) GetAllSongs(filters dto.GetSongsFilterDTO) ([]models.Song, error) {

	var AppendStrings []string
	queryString := "SELECT * FROM songs WHERE "

	// Добавляем фильтры
	if filters.Group != "" {
		appendString1 := "LOWER (group_name) LIKE LOWER ('%" + filters.Group + "%')"
		AppendStrings = append(AppendStrings, appendString1)
	}

	if filters.Song != "" {
		appendString2 := "LOWER (song_name) LIKE LOWER ('%" + filters.Song + "%')"
		AppendStrings = append(AppendStrings, appendString2)
	}

	if filters.ReleaseDate != "" {
		appendString3 := "release_date ='" + filters.ReleaseDate + "'"
		AppendStrings = append(AppendStrings, appendString3)
	}

	if filters.Text != "" {
		appendString4 := "LOWER (text) LIKE LOWER ('%" + filters.Text + "%')"
		AppendStrings = append(AppendStrings, appendString4)
	}

	if filters.Link != "" {
		appendString5 := "LOWER (link) LIKE LOWER ('%" + filters.Link + "%')"
		AppendStrings = append(AppendStrings, appendString5)
	}

	// Объединяем фильтры
	if len(AppendStrings) > 0 {
		queryString += strings.Join(AppendStrings, " AND ")
	}

	// Добавляем пагинацию
	offset := (filters.Page - 1) * filters.Limit
	queryString += fmt.Sprintf(" LIMIT %d OFFSET %d", filters.Limit, offset)

	log.Print(queryString)

	rows, err := s.DB.Query(queryString)
	if err != nil {
		return nil, fmt.Errorf("database query failed: %w", err)
	}
	defer rows.Close()

	var Songs []models.Song
	for rows.Next() {
		var song models.Song
		if err := rows.Scan(&song.Group, &song.Song, &song.Text, &song.ReleaseDate, &song.Link); err != nil {
			return nil, fmt.Errorf("error reading from database: %w", err)
		}

		Songs = append(Songs, song)
	}
	return Songs, nil
}
