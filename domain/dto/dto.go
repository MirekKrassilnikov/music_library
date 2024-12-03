package dto

// DTO для фильтров
type GetSongsFilterDTO struct {
	SongId      string `json:"song_id"`
	Group       string `json:"group"`
	Song        string `json:"song"`
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
	Page        string `json:"page"`
	Limit       string `json:"limit"`
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
	Page   string `json:"page"`
	Limit  string `json:"limit"`
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

type Pagination struct {
	CurrentPage int `json:"currentPage"`
	PageSize    int `json:"pageSize"`
	TotalItems  int `json:"totalItems"`
	TotalPages  int `json:"totalPages"`
	Offset      int `json:"offset"`
}
