package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/MirekKrassilnikov/music_library/domain/dto"
	"github.com/MirekKrassilnikov/music_library/domain/models"
)

type SongService struct {
	DB *sql.DB
}

// GetAllSongs обрабатывает бизнес-логику поиска песен
func (s *SongService) GetAllSongs(filters dto.GetSongsFilterDTO) ([]models.Song, error) {

	var AppendStrings []string
	queryString := "SELECT * FROM songs WHERE "

	if filters.Group != "" {
		appendString2 := "LOWER (group_name) LIKE LOWER ('%" + filters.Group + "%')"
		AppendStrings = append(AppendStrings, appendString2)
	}

	if filters.Song != "" {
		appendString3 := "LOWER (song_name) LIKE LOWER ('%" + filters.Song + "%')"
		AppendStrings = append(AppendStrings, appendString3)
	}

	if filters.ReleaseDate != "" {
		appendString4 := "release_date ='" + filters.ReleaseDate + "'"
		AppendStrings = append(AppendStrings, appendString4)
	}

	if filters.Text != "" {
		appendString5 := "LOWER (text) LIKE LOWER ('%" + filters.Text + "%')"
		AppendStrings = append(AppendStrings, appendString5)
	}

	if filters.Link != "" {
		appendString6 := "LOWER (link) LIKE LOWER ('%" + filters.Link + "%')"
		AppendStrings = append(AppendStrings, appendString6)
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

func (s *SongService) GetLyricsWithPagination(lyricsReq dto.GetLyricsDTO) ([]string, error) {

	text, err := s.GetLyricsById(lyricsReq.SongId)
	if err != nil {
		return nil, fmt.Errorf("failed to get lyrics: %w", err)
	}
	SplitedText := SplitIntoVerses(text)
	var end int
	if lyricsReq.Page > 0 {
		start := lyricsReq.Page - 1

		// Устанавливаем значение Limit (максимум 30)
		if lyricsReq.Limit > 30 {
			lyricsReq.Limit = 30
		} else if lyricsReq.Limit <= 0 { // Защита от некорректного значения Limit
			lyricsReq.Limit = 30 // Значение по умолчанию, если Limit не задан или <= 0
		}
		end = start + lyricsReq.Limit
		if end > len(SplitedText) {
			end = len(SplitedText)
		}
		if start > len(SplitedText) {
			start = len(SplitedText)
		}

		// Возвращаем нужный диапазон куплетов
		return SplitedText[start:end], nil
	}

	return SplitedText, nil
}

func (s *SongService) GetLyricsById(id string) (string, error) {
	queryString := "SELECT text FROM songs WHERE id = $1"

	log.Print(queryString)
	var text string
	err := s.DB.QueryRow(queryString, id).Scan(&text)
	if err != nil {
		return "", fmt.Errorf("error reading from database: %w", err)
	}
	return text, nil
}

func SplitIntoVerses(text string) []string {
	verses := strings.Split(text, "\n\n")
	for i, verse := range verses {
		verses[i] = strings.TrimSpace(verse)
	}

	return verses
}

func (s *SongService) DeleteSong(id string) (string, error) {
	queryString := "DELETE FROM songs WHERE id = $1"
	_, err := s.DB.Exec(queryString, id)
	if err != nil {
		return "", fmt.Errorf("error deleting song from database: %w", err)
	}

	return fmt.Sprintf("Песня с id %s успешно удалена", id), nil
}

func (s *SongService) AddNewSong(group, song string) (string, error) {
	addictionalInfo, err := getAddictionalInfo(group, song)
	if err != nil {
		return "", fmt.Errorf("failed to get addictnional info: %w", err)
	}
	newSong := dto.SongDTO{
		Group:       group,
		Song:        song,
		Text:        addictionalInfo.Text,
		ReleaseDate: addictionalInfo.ReleaseDate,
		Link:        addictionalInfo.Link,
	}

	query := `
		INSERT INTO songs (group_name, song_name, text, release_date, link) 
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`

	var songID string
	err = s.DB.QueryRow(query, newSong.Group, newSong.Song, newSong.Text, newSong.ReleaseDate, newSong.Link).Scan(&songID)
	if err != nil {
		return "", fmt.Errorf("failed to insert song into database: %w", err)
	}

	// Возвращаем ID новой песни
	return songID, nil
}

func getAddictionalInfo(group, song string) (dto.AddictionalInfo, error) {
	// Формируем URL для запроса
	apiUrl := fmt.Sprintf("http://localhost:8080/info?group=%s&song=%s", group, song)

	// Выполняем HTTP-запрос
	resp, err := http.Get(apiUrl)
	if err != nil {
		return dto.AddictionalInfo{}, fmt.Errorf("failed to call external API: %w", err)
	}
	defer resp.Body.Close()
	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		return dto.AddictionalInfo{}, fmt.Errorf("API request failed with status: %s", resp.Status)
	}
	// Читаем тело ответа
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return dto.AddictionalInfo{}, fmt.Errorf("failed to read response body: %w", err)
	}

	var addictionalInfo dto.AddictionalInfo
	if err := json.Unmarshal(body, &addictionalInfo); err != nil {
		return dto.AddictionalInfo{}, fmt.Errorf("failed to parse JSON: %w", err)
	}
	return addictionalInfo, nil
}
