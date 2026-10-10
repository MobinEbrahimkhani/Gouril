package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"
	"net/url"

	"github.com/MobinEbrahimkhani/Gouril/database"
)

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	ShortCode string `json:"short_code"`
	ShortURL  string `json:"short_url"`
}

func generateShortCode() string {
	const characters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	code := make([]byte, 6)

	for i := range code {
		code[i] = characters[rand.Intn(len(characters))]
	}

	return string(code)
}

func ShortenHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var request ShortenRequest

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if request.URL == "" {
			http.Error(w, "URL is required", http.StatusBadRequest)
			return
		}

		parsedURL, err := url.ParseRequestURI(request.URL)
		if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
			http.Error(w, "Invalid URL", http.StatusBadRequest)
			return
		}

		var shortCode string

		for {
			shortCode = generateShortCode()

			exists, err := database.ShortCodeExists(db, shortCode)
			if err != nil {
				http.Error(w, "Database error", http.StatusInternalServerError)
				return
			}

			if !exists {
				break
			}
		}

		if err := database.InsertURL(db, shortCode, request.URL); err != nil {
			http.Error(w, "Could not save URL", http.StatusInternalServerError)
			return
		}

		response := ShortenResponse{
			ShortCode: shortCode,
			ShortURL:  "http://localhost:8080/" + shortCode,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		if err := json.NewEncoder(w).Encode(response); err != nil {
			return
		}
	}
}

func RedirectHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if r.URL.Path == "/" {
			http.Error(w, "Short code is required", http.StatusBadRequest)
			return
		}

		shortCode := r.URL.Path[1:]

		originalURL, err := database.GetURL(db, shortCode)
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Short URL not found", http.StatusNotFound)
			return
		}
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, originalURL, http.StatusFound)
	}
}
