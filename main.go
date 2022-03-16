package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/blevesearch/bleve/v2"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

type Event struct {
	ID          int
	Name        string
	Description string
	Local       string
	Website     string
	Start       time.Time
	End         time.Time
}

const (
	testIdx = "test.bleve"
	dbFile  = "test.sqlite3.db"
)

func idxCreate() bleve.Index {
	idx, _ := CreateIndex(testIdx)
	return idx
}

// create a SQLite3 database file, create an events table and fill with some data.
func dbCreate() (*gorm.DB, []Event) {
	db, err := gorm.Open("sqlite3", dbFile)
	if err != nil {
		fmt.Println("ERROR: sqlite3 open", err)
	}
	db.DropTableIfExists(&Event{})
	db.CreateTable(&Event{})
	eventList := fillDatabase(db)
	return db, eventList
}

// fill the database with some data
func fillDatabase(db *gorm.DB) []Event {
	eventList := EventList
	// inserting the events
	for _, event := range eventList {
		db.Create(&event)
	}

	return eventList
}

func CreateIndex(path string) (bleve.Index, error) {
	mapping := bleve.NewIndexMapping()
	index, err := bleve.New(path, mapping)
	if err != nil {
		return nil, err
	}
	return index, nil
}

func (e *Event) Index(index bleve.Index) error {
	err := index.Index(strconv.Itoa(e.ID), e)
	return err
}
func indexEvents(idx bleve.Index, eventList []Event) {
	for _, event := range eventList {
		event.Index(idx)
	}
}

func SearchForQuery(searchTerm string) {
	db, eventList := dbCreate()
	idx := idxCreate()
	indexEvents(idx, eventList)
	query := bleve.NewMatchQuery(searchTerm)
	searchRequest := bleve.NewSearchRequest(query)
	searchResults, err := idx.Search(searchRequest)
	if err != nil {
		fmt.Println("Search ERROR:", err)
		return
	}
	event := &Event{}
	db.First(event, searchResults.Hits[0].ID)
	fmt.Println("SEARCH RESULTS: ", searchResults.Status.Successful, searchResults.Total)
	if event.Name != "dotGo 2015" {
		fmt.Println("Expected \"dotGo 2015\", Receive: ", event.Name)
	} else {
		fmt.Println("Should return an event with the name equal a", event.Name)
	}

	defer idxDestroy()
	defer dbDestroy()
}

func main() {
	// SearchForQuery("http://www.gophercon.in/")
	createDB("blog.db")
}

func idxDestroy() {
	os.RemoveAll(testIdx)
}

func dbDestroy() {
	os.RemoveAll(dbFile)
}
