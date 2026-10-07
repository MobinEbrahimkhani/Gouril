package handler

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"net/url"
)

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	ShortCode string `json:"short_code"`
	ShortURL  string `json:"short_url"`
}

type URL struct {
	OriginalURL string
}

var URLs = make(map[string]URL)

func generateShortCode() string {
	const characters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	code := make([]byte, 6)

	for i := range code {
		code[i] = characters[rand.Intn(len(characters))]
	}

	return string(code)
}

func ShortenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request ShortenRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
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

		if _, exists := URLs[shortCode]; !exists {
			break
		}
	}

	URLs[shortCode] = URL{
		OriginalURL: request.URL,
	}

	response := ShortenResponse{
		ShortCode: shortCode,
		ShortURL:  "http://localhost:8080/" + shortCode,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(response)
}


func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path == "/" {
		http.Error(w, "Short code is required", http.StatusBadRequest)
		return
	}

	shortCode := r.URL.Path[1:]

	url, exists := URLs[shortCode]
	if !exists {
		http.Error(w, "Short URL not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, url.OriginalURL, http.StatusFound)
}
