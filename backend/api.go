package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"
)

type EntryRequest struct {
	VisitedAt time.Time `json:"visited_at"`
	Place     string    `json:"place"`
	Comment   string    `json:"comment"`
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func apiEntriesHandler(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		entries, err := store.ListEntries()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		err = json.NewEncoder(w).Encode(entries)
		if err != nil {
			http.Error(w, "JSON encoding failed", http.StatusInternalServerError)
			return
		}
	}
}

func apiSingleEntryHandler(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		entry, err := store.GetEntry(id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w, "Entry not found", http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		err = json.NewEncoder(w).Encode(entry)
		if err != nil {
			http.Error(w, "JSON encoding failed", http.StatusInternalServerError)
			return
		}
	}
}

func apiDeleteHandler(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		err = store.DeleteEntry(id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w, "Entry not found", http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func apiCreateEntryHandler(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req EntryRequest

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		entry, err := store.AddEntry(req.VisitedAt, req.Place, req.Comment)
		if err != nil {
			http.Error(w, "Failed to create entry", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		err = json.NewEncoder(w).Encode(entry)
		if err != nil {
			http.Error(w, "JSON encoding failed", http.StatusInternalServerError)
			return
		}
	}
}

func apiEditHandler(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			http.Error(w, "Invalid id", http.StatusBadRequest)
			return
		}

		var req EntryRequest

		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		entry, err := store.UpdateEntry(id, req.VisitedAt, req.Place, req.Comment)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w, "Entry not found", http.StatusNotFound)
				return
			}

			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(entry)
		if err != nil {
			http.Error(w, "JSON encoding failed", http.StatusInternalServerError)
			return
		}
	}
}
