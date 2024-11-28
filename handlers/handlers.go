package handlers

import (
	"net/http"
	"strconv"

	"github.com/MirekKrassilnikov/music_library/dto"
	"github.com/MirekKrassilnikov/music_library/services"
)



func HandleGetAllSongs(w http.ResponseWriter, r *http.Request) {
	// создаем сервис
	songService := services.SongService {
		DB: db
	}

	var page int
	var limit int
	if r.URL.Query().Get("page") != "" {
		var err error
		page, err := strconv.Atoi(r.URL.Query().Get("page"))
		if err != nil || page < 1 {
			page = 1
		}
	}
	if r.URL.Query().Get("limit") != "" {
		var err error
		limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil || limit < 1 {
			limit = 10
		}
	}

	GetSongsFilterDTO := dto.GetSongsFilterDTO{
		Group:       r.URL.Query().Get("group"),
		Song:        r.URL.Query().Get("song"),
		ReleaseDate: r.URL.Query().Get("releaseDate"),
		Text:        r.URL.Query().Get("text"),
		Link:        r.URL.Query().Get("link"),
		Page:        page,
		Limit:       limit,
	}

	getAllSongsResponse := songService.GetAllSongs(GetSongsFilterDTO)
}
