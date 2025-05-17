package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"time"
)

func main() {
	srv := &http.Server{
		Addr:         "8080",
		Handler:      mux(),
		ReadTimeout:  2 * time.Minute,
		WriteTimeout: 2 * time.Minute,
	}
	log.Fatal(http.ListenAndServe(":8080", srv.Handler))
}

func mux() *chi.Mux {
	r := chi.NewRouter()

	// Add CORS middleware
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"}, // or restrict to specific domains
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Max value = 600 (in seconds)
	}))

	// Your routes
	r.Route("/hello", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("helloooo"))
		})
	})

	r.Post("/manga", func(w http.ResponseWriter, r *http.Request) {
		type Payload struct {
			MinChapter string `json:"minChapter"`
			MaxChapter string `json:"maxChapter"`
			MangaURL   string `json:"url"`
		}

		var p Payload
		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// Optional: log values for debugging
		log.Println("Received:", p.MinChapter, p.MaxChapter, p.MangaURL)
		Manga(p.MinChapter, p.MaxChapter, p.MangaURL)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK) // optional; 200 is the default
		w.Write([]byte(`{"message": "Downloaded successfully!"}`))
	})

	return r
}
