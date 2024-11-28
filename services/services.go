package services

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/MirekKrassilnikov/music_library/models"
)

func GetAllSongs(w http.ResponseWriter, r *http.Request) {
	group := r.URL.Query().Get("group")
	song := r.URL.Query().Get("song")
	releaseDate := r.URL.Query().Get("releaseDate")
	text := r.URL.Query().Get("text")
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

	var AppendStrings []string
	queryString := "SELECT * FROM songs WHERE "

	if group != "" {
		appendString1 := "LOWER (group_name) LIKE LOWER ('%" + group + "%')"
		AppendStrings = append(AppendStrings, appendString1)
	}

	if song != "" {
		appendString2 := "LOWER (song_name) LIKE LOWER ('%" + song + "%')"
		AppendStrings = append(AppendStrings, appendString2)
	}

	if releaseDate != "" {
		appendString3 := "release_date ='" + releaseDate + "'"
		AppendStrings = append(AppendStrings, appendString3)
	}

	if text != "" {
		appendString4 := "LOWER (text) LIKE LOWER ('%" + text + "%')"
		AppendStrings = append(AppendStrings, appendString4)
	}

	if link != "" {
		appendString5 := "LOWER (link) LIKE LOWER ('%" + link + "%')"
		AppendStrings = append(AppendStrings, appendString5)
	}

	if len(AppendStrings) > 0 {
		queryString += strings.Join(AppendStrings, " AND ")
	}

	log.Print(queryString)

	rows, err := db.Query(queryString)
	if err != nil {
		http.Error(w, "Database query failed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var Songs []models.Song
	for rows.Next() {
		var song models.Song
		if err := rows.Scan(&song.Group, &Song.Song, &song.Text, &song.ReleaseDate, &song.Link); err != nil {
			http.Error(w, "Error reading from database", http.StatusInternalServerError)
			return
		}
		Songs = append(Songs, song)
	}

}
