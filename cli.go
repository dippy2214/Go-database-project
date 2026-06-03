package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

// --------------------------------------------------------
//					CLI TOOLING
//---------------------------------------------------------

// this cli tool gives some useful testing functionality from
// the terminal, as an easy way to interact with and test the
// database functionality

func displayEntries(entries []Entry) {

	for _, e := range entries {
		fmt.Printf("%d | %s | %s | %s\n", e.ID, e.VisitedAt.Format("2006-01-02 15:04:05"), e.Place, e.Comment)
	}
}

func list(store *Store) {
	e, err := store.ListEntries()
	if err != nil {
		log.Fatal("failed to fetch list: ", err)
	}
	displayEntries(e)
}

func add(store *Store) {
	if len(os.Args) < 5 {
		fmt.Println("usage: add <time (YYYY-MM-DD HH:MM:SS)> <place> <comment>")
		return
	}
	t, err := time.Parse("2006-01-02 15:04:05", os.Args[2])
	if err != nil {
		log.Fatal("failed to parse time: ", err)
	}
	err = store.AddEntry(t, os.Args[3], os.Args[4])
	if err != nil {
		log.Fatal("failed to add entry: ", err)
	}
}

func recent(store *Store) {
	if len(os.Args) < 3 {
		fmt.Println("usage: recent <count>")
		return
	}
	count, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Fatal("failed to parse count: ", err)
	}
	e, err := store.RecentEntries(count)
	if err != nil {
		log.Fatal("failed to fetch list: ", err)
	}
	displayEntries(e)
}
