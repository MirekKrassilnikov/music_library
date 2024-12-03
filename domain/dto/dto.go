package dto

// DTO для фильтров
type GetSongsFilterDTO struct {
	SongId      string `json:"song_id"`
	Group       string `json:"group"`
	Song        string `json:"song"`
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
	Page        int    `json:"page"`
	Limit       int    `json:"limit"`
}

// DTO для ответа
type SongDTO struct {
	Group       string `json:"group"`
	Song        string `json:"song"`
	Text        string `json:"text"`
	ReleaseDate string `json:"release_date"`
	Link        string `json:"link"`
}

type GetLyricsDTO struct {
	SongId string `json:"song_id"`
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
}

type NewSongDTO struct {
	Group string `json:"group"`
	Song  string `json:"song"`
}

type AddictionalInfo struct {
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}
