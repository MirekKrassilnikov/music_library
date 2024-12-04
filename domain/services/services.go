package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/MirekKrassilnikov/music_library/domain/dto"
	"github.com/MirekKrassilnikov/music_library/domain/models"
)

type SongService struct {
	DB *sql.DB
}

func (s *SongService) GetAllSongs(filters dto.GetSongsFilterDTO) ([]models.Song, dto.Pagination, error) {

	var conditions []string
	var args []interface{}

	// Формируем условия фильтрации и добавляем параметры
	if filters.Group != "" {
		conditions = append(conditions, "LOWER(group_name) LIKE LOWER(?)")
		args = append(args, "%"+filters.Group+"%")
	}

	if filters.Song != "" {
		conditions = append(conditions, "LOWER(song_name) LIKE LOWER(?)")
		args = append(args, "%"+filters.Song+"%")
	}

	if filters.ReleaseDate != "" {
		conditions = append(conditions, "release_date = ?")
		args = append(args, filters.ReleaseDate)
	}

	if filters.Text != "" {
		conditions = append(conditions, "LOWER(text) LIKE LOWER(?)")
		args = append(args, "%"+filters.Text+"%")
	}

	if filters.Link != "" {
		conditions = append(conditions, "LOWER(link) LIKE LOWER(?)")
		args = append(args, "%"+filters.Link+"%")
	}

	// Создаем часть запроса с WHERE
	queryString := "SELECT * FROM songs"
	if len(conditions) > 0 {
		queryString += " WHERE " + strings.Join(conditions, " AND ")
	}

	// Подсчет количества записей для пагинации
	countQuery := "SELECT COUNT(*) FROM songs"
	if len(conditions) > 0 {
		countQuery += " WHERE " + strings.Join(conditions, " AND ")
	}
	var totalCount int
	err := s.DB.QueryRow(countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, dto.Pagination{}, fmt.Errorf("database query failed: %w", err)
	}

	// Получаем информацию о пагинации
	Pagination, err := s.GetPagination(totalCount, filters.Page, filters.Limit)
	if err != nil {
		return nil, dto.Pagination{}, fmt.Errorf("pagination calculation failed: %w", err)
	}

	// Добавляем LIMIT и OFFSET к запросу
	queryString += fmt.Sprintf(" LIMIT ? OFFSET ?")
	args = append(args, filters.Limit, Pagination.Offset)

	// Выполняем запрос
	rows, err := s.DB.Query(queryString, args...)
	if err != nil {
		return nil, dto.Pagination{}, fmt.Errorf("database query failed: %w", err)
	}
	defer rows.Close()

	var Songs []models.Song
	for rows.Next() {
		var song models.Song
		if err := rows.Scan(&song.Group, &song.Song, &song.Text, &song.ReleaseDate, &song.Link); err != nil {
			return nil, dto.Pagination{}, fmt.Errorf("error reading from database: %w", err)
		}

		Songs = append(Songs, song)
	}

	return Songs, Pagination, nil
}

func (s *SongService) GetPagination(totalCount int, pageStr string, limitStr string) (dto.Pagination, error) {

	// Parsing page number
	page := 1
	if pageStr != "" {
		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			page = 1
		}
	}

	// Parsing Limit
	limit := 10
	if limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit < 1 {
			limit = 10
		}
	}

	// Calculation total ammount of pages
	totalPages := (totalCount + limit - 1) / limit

	// Calculating offset
	offset := (page - 1) * limit

	pagination := dto.Pagination{
		CurrentPage: page,
		PageSize:    limit,
		TotalItems:  totalCount,
		TotalPages:  totalPages,
		Offset:      offset, // Добавляем offset, если нужен для запроса
	}

	return pagination, nil
}

func (s *SongService) GetLyricsWithPagination(lyricsReq dto.GetLyricsDTO) ([]string, dto.Pagination, error) {

	text, err := s.GetLyricsById(lyricsReq.SongId)
	if err != nil {
		return nil, dto.Pagination{}, fmt.Errorf("failed to get lyrics: %w", err)
	}
	SplitedText := SplitIntoVerses(text)
	totalCount := len(SplitedText)

	if lyricsReq.Page != "" {
		Pagination, err := s.GetPagination(totalCount, lyricsReq.Page, lyricsReq.Limit)
		if err != nil {
			return nil, dto.Pagination{}, fmt.Errorf("pagination calculation failed: %w", err)
		}

		end := Pagination.Offset + Pagination.PageSize

		// Returning requested range
		return SplitedText[Pagination.Offset:end], Pagination, nil
	}

	return SplitedText, dto.Pagination{}, nil
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

	apiUrl := os.Getenv("API_URL")
	if apiUrl == "" {
		return dto.AddictionalInfo{}, fmt.Errorf("API_URL environment variable is not set")
	}

	url := fmt.Sprintf("%s/info?group=%s&song=%s", apiUrl, group, song)

	// Sending request
	resp, err := http.Get(url)
	if err != nil {
		return dto.AddictionalInfo{}, fmt.Errorf("failed to call external API: %w", err)
	}
	defer resp.Body.Close()
	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		return dto.AddictionalInfo{}, fmt.Errorf("API request failed with status: %s", resp.Status)
	}
	// Reading body
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
