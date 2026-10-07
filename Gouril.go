package main

import (
	"fmt"
	"net/http"

	"github.com/MobinEbrahimkhani/Gouril/handler"

	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	http.HandleFunc("/shorten", handler.ShortenHandler)
	http.HandleFunc("/", handler.RedirectHandler)

	fmt.Println("Gouril is listening on http://localhost:8080")

	server := &http.Server{
		Addr: ":8080",
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-stop

		fmt.Println("\nShutting down server...")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			fmt.Println("Server shutdown error:", err)
		}
	}()

	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		fmt.Println(err)
	}
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path == "/" {
		http.Error(w, "Short code is required", http.StatusBadRequest)
		return
	}

	shortCode := r.URL.Path[1:]

	url, exists := handler.URLs[shortCode]
	if !exists {
		http.Error(w, "Short URL not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, url.OriginalURL, http.StatusFound)
}
