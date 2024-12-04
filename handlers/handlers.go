package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/MirekKrassilnikov/music_library/domain/dto"
	"github.com/MirekKrassilnikov/music_library/domain/services"
)

type SongHandler struct {
	SongService *services.SongService
	Logger      *slog.Logger
}

// Handler for GET /songs - returns list of songs
func (h *SongHandler) HandleGetAllSongs(w http.ResponseWriter, r *http.Request) {

	GetSongsFilterDTO := dto.GetSongsFilterDTO{
		SongId:      r.URL.Query().Get("id"),
		Group:       r.URL.Query().Get("group"),
		Song:        r.URL.Query().Get("song"),
		ReleaseDate: r.URL.Query().Get("releaseDate"),
		Text:        r.URL.Query().Get("text"),
		Link:        r.URL.Query().Get("link"),
		Page:        r.URL.Query().Get("page"),
		Limit:       r.URL.Query().Get("limit"),
	}
	h.Logger.Debug("Filter parameters", slog.Any("filters", GetSongsFilterDTO))

	songs, Pagination, err := h.SongService.GetAllSongs(GetSongsFilterDTO)
	if err != nil {
		http.Error(w, "Failed to get songs: "+err.Error(), http.StatusInternalServerError)
		return
	}
	response := map[string]interface{}{
		"songs":      songs,
		"pagination": Pagination,
	}

	responseData, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	h.Logger.Info("Successfully retrieved songs", slog.Any("response", response))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseData)
}

func (h *SongHandler) HandleGetLyrics(w http.ResponseWriter, r *http.Request) {

	GetLyricsDTO := dto.GetLyricsDTO{
		SongId: r.URL.Query().Get("id"),
		Page:   r.URL.Query().Get("page"),
		Limit:  r.URL.Query().Get("limit"),
	}
	h.Logger.Debug("Filter parameters", slog.Any("filters", GetLyricsDTO))

	lyrics, Pagination, err := h.SongService.GetLyricsWithPagination(GetLyricsDTO)
	if err != nil {
		http.Error(w, "Failed to get lyrics: "+err.Error(), http.StatusInternalServerError)
		return
	}
	response := map[string]interface{}{
		"lyrics":     lyrics,
		"pagination": Pagination,
	}

	responseData, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	h.Logger.Info("Successfully retrieved lyrics", slog.Any("response", response))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseData)
}

func (h SongHandler) HandleDeleteSong(w http.ResponseWriter, r *http.Request) {
	songId := r.URL.Query().Get("id")
	h.Logger.Debug("Delete parameters", slog.String("songId", songId))

	response, err := h.SongService.DeleteSong(songId)
	if err != nil {
		h.Logger.Error("Failed to delete song", slog.String("error", err.Error()))
		http.Error(w, "Failed to delete song: "+err.Error(), http.StatusInternalServerError)
		return
	}

	responseData, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	h.Logger.Info("Successfully deleted song", slog.String("songId", songId))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseData)
}

func (h SongHandler) HandleAddNewSong(w http.ResponseWriter, r *http.Request) {

	NewSong := dto.NewSongDTO{
		Group: r.URL.Query().Get("group"),
		Song:  r.URL.Query().Get("song"),
	}

	h.Logger.Debug("New song parameters", slog.Any("newSong", NewSong))

	returnedId, err := h.SongService.AddNewSong(NewSong.Group, NewSong.Song)
	if err != nil {
		http.Error(w, "Failed to add song: "+err.Error(), http.StatusInternalServerError)
		return
	}

	h.Logger.Info("Successfully added new song", slog.String("songId", returnedId))
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
