package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/MirekKrassilnikov/music_library/domain/dto"
	"github.com/MirekKrassilnikov/music_library/domain/services"
)

type SongHandler struct {
	SongService *services.SongService
}

// Handler for GET /songs - returns list of songs
func (h *SongHandler) HandleGetAllSongs(w http.ResponseWriter, r *http.Request) {

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
		SongId:      r.URL.Query().Get("id"),
		Group:       r.URL.Query().Get("group"),
		Song:        r.URL.Query().Get("song"),
		ReleaseDate: r.URL.Query().Get("releaseDate"),
		Text:        r.URL.Query().Get("text"),
		Link:        r.URL.Query().Get("link"),
		Page:        page,
		Limit:       limit,
	}

	getAllSongsResponse, err := h.SongService.GetAllSongs(GetSongsFilterDTO)
	if err != nil {
		http.Error(w, "Failed to get songs: "+err.Error(), http.StatusInternalServerError)
		return
	}

	responseData, err := json.Marshal(getAllSongsResponse)
	if err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseData)
}

func (h *SongHandler) HandleGetLyrics(w http.ResponseWriter, r *http.Request) {
	// Переводим параметры в int
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

	GetLyricsDTO := dto.GetLyricsDTO{
		SongId: r.URL.Query().Get("id"),
		Page:   page,
		Limit:  limit,
	}

	lyrics, err := h.SongService.GetLyricsWithPagination(GetLyricsDTO)
	if err != nil {
		http.Error(w, "Failed to get lyrics: "+err.Error(), http.StatusInternalServerError)
		return
	}
	responseData, err := json.Marshal(lyrics)
	if err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseData)
}

func (h SongHandler) HandleDeleteSong(w http.ResponseWriter, r *http.Request) {
	songId := r.URL.Query().Get("id")

	response, err := h.SongService.DeleteSong(songId)
	if err != nil {
		http.Error(w, "Failed to delete song: "+err.Error(), http.StatusInternalServerError)
		return
	}

	responseData, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseData)
}

func (h SongHandler) HandleAddNewSong(w http.ResponseWriter, r *http.Request) {

	NewSong := dto.NewSongDTO{
		Group: r.URL.Query().Get("group"),
		Song:  r.URL.Query().Get("song"),
	}

	returnedId, err := h.SongService.AddNewSong(NewSong.Group, NewSong.Song)
	if err != nil {
		http.Error(w, "Failed to add song: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := "Succesfully added song with id: " + returnedId
	responseData, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseData)

}
