package dto

type SongDTO struct {
	Group       string `json:"group"`
	Song        string `json:"song"`
	Text        string `json:"text"`
	ReleaseDate string `json:"release_date"`
	Link        string `json:"link"`
}
