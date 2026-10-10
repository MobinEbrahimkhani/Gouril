package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MobinEbrahimkhani/Gouril/database"
	"github.com/MobinEbrahimkhani/Gouril/handler"
)

func main() {
	db, err := database.Open()
	if err != nil {
		fmt.Println("Database error:", err)
		return
	}
	defer db.Close()

	if err := database.Init(db); err != nil {
		fmt.Println("Database initialization error:", err)
		return
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/shorten", handler.ShortenHandler(db))
	mux.HandleFunc("/", handler.RedirectHandler(db))

	server := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(stop)

	go func() {
		<-stop

		fmt.Println("\nShutting down server...")

		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			5*time.Second,
		)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			fmt.Println("Server shutdown error:", err)
		}
	}()

	fmt.Println("Database connected successfully")
	fmt.Println("Gouril is listening on http://localhost:8080")

	err = server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		fmt.Println("Server error:", err)
	}
}
