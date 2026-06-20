package main

import (
	"html/template"
	"net/http"
	"strconv"
	"time"
)

// this http part of the project deliberately separates the view data
// from the stored data. This means that data is easily reformatted
// to suit a different time and date format in this example

type EntryView struct {
	ID        int
	VisitedAt string
	Place     string
	Comment   string
}

func entriesHandler(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		entries, err := store.ListEntries()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		tmpl, err := template.ParseFiles(("templates/entries.html"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		views := make([]EntryView, 0, len(entries))

		for _, e := range entries {
			views = append(views, EntryView{
				ID:        e.ID,
				VisitedAt: e.VisitedAt.Format("2006-01-02 15:04"),
				Place:     e.Place,
				Comment:   e.Comment,
			})
		}

		err = tmpl.Execute(w, views)
	}
}

func addFormHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/add.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}

func addPostHandler(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		err := r.ParseForm()
		if err != nil {
			http.Error(w, "bad form data", http.StatusBadRequest)
			return
		}

		timeStr := r.FormValue("visited_at")
		place := r.FormValue("place")
		comment := r.FormValue("comment")

		t, err := time.Parse("2006-01-02 15:04:05", timeStr)
		if err != nil {
			http.Error(w, "invalid time format", http.StatusBadRequest)
			return
		}

		_, err = store.AddEntry(t, place, comment)
		if err != nil {
			http.Error(w, "db insert failed", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/entries", http.StatusSeeOther)
	}
}

func deleteHandler(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		r.ParseForm()

		idStr := r.FormValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		err = store.DeleteEntry(id)
		if err != nil {
			http.Error(w, "delete failed", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/entries", http.StatusSeeOther)
	}
}

func editGetHandler(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		r.ParseForm()

		idStr := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		entry, err := store.GetEntry(id)
		if err != nil {
			http.Error(w, "delete failed", http.StatusInternalServerError)
			return
		}

		tmpl, err := template.ParseFiles("templates/edit.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var view EntryView
		view.ID = entry.ID
		view.VisitedAt = entry.VisitedAt.Format("2006-01-02 15:04:05")
		view.Place = entry.Place
		view.Comment = entry.Comment

		tmpl.Execute(w, view)
	}
}

func editPostHandler(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		r.ParseForm()

		id, _ := strconv.Atoi(r.FormValue("id"))
		place := r.FormValue("place")
		comment := r.FormValue("comment")

		t, err := time.Parse("2006-01-02 15:04:05", r.FormValue("visited_at"))
		if err != nil {
			http.Error(w, "invalid time", http.StatusBadRequest)
			return
		}

		_, err = store.UpdateEntry(id, t, place, comment)
		if err != nil {
			http.Error(w, "update failed", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/entries", http.StatusSeeOther)
	}
}

func homeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/entries", http.StatusSeeOther)
	}
}
