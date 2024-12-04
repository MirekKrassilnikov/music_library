package dto

type GetLyricsDTO struct {
	SongId string `json:"song_id"`
	Page   string `json:"page"`
	Limit  string `json:"limit"`
}
