package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// here in main we handle running the program and parsing arguments
// this allows for the cli tool and web hosting to coexist in the same program

func main() {

	if len(os.Args) < 2 {
		fmt.Println("expected 'list' or 'add'")
		return
	}

	store := mustOpenStore()

	switch os.Args[1] {
	case "list":
		list(store)
	case "add":
		add(store)
	case "recent":
		recent(store)
	case "serve":
		http.HandleFunc("/entries", entriesHandler(store))
		http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				addFormHandler(w, r)
			case http.MethodPost:
				addPostHandler(store)(w, r)
			default:
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			}
		})
		http.HandleFunc("/delete", deleteHandler(store))
		http.HandleFunc("/edit", func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				editGetHandler(store)(w, r)
			case http.MethodPost:
				editPostHandler(store)(w, r)
			default:
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			}
		})
		log.Fatal(http.ListenAndServe(":8080", nil))
	default:
		fmt.Println("unknown command: ", os.Args[1])
	}
}
