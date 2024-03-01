package services

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"linkSwitch/database"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func RunServer(db database.Database) http.Handler {
	router := chi.NewRouter()

	router.Route("/", func(r chi.Router) {
		r.Post("/", createURL(db))
		r.Get("/{shortURL}", getURL(db))
		r.Get("/stats", getStats(db))
		r.Get("/stats/{shortURL}", getClicks(db))
	})

	return router
}

func getClicks(db database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		shortURL := chi.URLParam(r, "shortURL")

		url, err := db.FindOne(shortURL)
		if err != nil {
			http.Error(w, "url not found", http.StatusNotFound)
			log.Println("url not found:", err)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(strconv.Itoa(url.Clicks)))
	}
}

func createURL(db database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var url database.URL
		err := json.NewDecoder(r.Body).Decode(&url)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}

		h := fnv.New32a()
		h.Write([]byte(url.Long))
		hash := h.Sum32()

		url.Short = fmt.Sprintf("%x", hash)[:8]

		err = db.InsertOne(url)
		if err != nil {
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			log.Println("Failed to create user:", err)
			return
		}

		jsonData, err := json.Marshal(url)
		if err != nil {
			http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		_, err = w.Write(jsonData)
		if err != nil {
			log.Println("Failed to write response:", err)
		}
	}
}

func getURL(db database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		shortURL := chi.URLParam(r, "shortURL")

		url, err := db.FindOne(shortURL)
		if err != nil {
			http.Error(w, "url not found", http.StatusNotFound)
			log.Println("url not found:", err)
			return
		}

		url.Clicks += 1

		url, err = db.Update(*url)
		if err != nil {
			http.Error(w, "update failed", http.StatusBadRequest)
			log.Println("update failed", err)
			return
		}

		http.Redirect(w, r, url.Long, http.StatusFound)
	}
}

func getStats(db database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		urls, err := db.Find()
		if err != nil {
			http.Error(w, "Failed to get users", http.StatusInternalServerError)
			log.Println("Failed to get users:", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(strconv.Itoa(len(urls))))
	}
}
