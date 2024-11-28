package services

import (
	"net/http"
	"strconv"
)

func GetAllSongs(w http.ResponseWriter, r *http.Request) {
	group := r.URL.Query().Get("group")
	song := r.URL.Query().Get("song")
	text := r.URL.Query().Get("text")
	releaseDate := r.URL.Query().Get("releaseDate")
	link := r.URL.Query().Get("link")

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	if pageStr != "" {
		var err error
		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			page = 1
		}
	}
	if limitStr != "" {
		var err error
		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit < 1 {
			limit = 10
		}
	}

}
