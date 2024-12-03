package routes

import (
	"net/http"

	"github.com/MirekKrassilnikov/music_library/handlers"
)

func RegisterRoutes(mux *http.ServeMux, songHandler *handlers.SongHandler) {
	mux.HandleFunc("/songs", songHandler.HandleGetAllSongs)
	mux.HandleFunc("/lyrics", songHandler.HandleGetLyrics)
	mux.HandleFunc("/delete", songHandler.HandleDeleteSong)
	mux.HandleFunc("/add", songHandler.HandleAddNewSong)
}
