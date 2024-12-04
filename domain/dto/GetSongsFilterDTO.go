package dto

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
