package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"time"

	"github.com/boltdb/bolt"
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

var db *bolt.DB
var err, anerr error
var bucketName = []byte("evtBucket")
var fieldName = []byte("eventsKey3")

func OpenDB() {
	db, err = bolt.Open("event.db", 0600, nil)
	if err != nil {
		fmt.Println("Open ERROR:", err)
	}
	fmt.Println(db.GoString(), db.Info())
}

func CreateEvent(events *Event) {
	anerr = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(bucketName)
		if err != nil {
			fmt.Println("CreateBucketIfNotExists: ", err)
		}

		// Generate ID for the events.
		id, _ := b.NextSequence()
		events.ID = int(id)
		fmt.Println(id, events.ID)

		// Marshal events data into bytes.
		buf, err := json.Marshal(events)
		if err != nil {
			fmt.Println("Marshall ERROR: ", err)
		}

		// Persist bytes to users bucket.
		err = b.Put(itob(events.ID), buf)

		// //For now let's concentrate on just static ID index
		// err = b.Put(fieldName, buf)

		// Additional Get for checking the field added
		v := b.Get(itob(events.ID))
		fmt.Printf("The answer is: %s\n", v)
		return err
	})
	if anerr != nil {
		fmt.Println("Update Error:", anerr)
	}
}

// itob returns an 8-byte big endian representation of v.
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

// btoi converts from byte to int
func btoi(b []byte) int {
	return int(binary.BigEndian.Uint64(b))
}

func ReadDB() {
	anerr = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		v := b.Get(fieldName)
		fmt.Printf("The answer is: %s\n", v)
		return nil
	})
	if anerr != nil {
		fmt.Println("View ERROR:", anerr)
	}
}
func ReadDBStream() {
	db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket(bucketName)

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			fmt.Printf("key=%d, value=%s\n", btoi(k), v)
		}

		return nil
	})
}
func anothermain() {
	// events := []Event{
	// 	{1, "dotGo 2015", "The European Go conference", "Paris", "http://www.dotgo.eu/", time.Date(2015, 11, 19, 9, 0, 0, 0, time.UTC), time.Date(2015, 11, 19, 18, 30, 0, 0, time.UTC)},

	// 	{2, "GopherCon INDIA 2016", "The Go Conference in India", "Bengaluru", "http://www.gophercon.in/", time.Date(2016, 2, 19, 0, 0, 0, 0, time.UTC), time.Date(2016, 2, 20, 23, 59, 0, 0, time.UTC)},

	// 	{3, "GopherCon 2016", "GopherCon, It is the largest event in the world dedicated solely to the Go programming language. It's attended by the best and the brightest of the Go team and community.", "Denver", "http://gophercon.com/", time.Date(2016, 7, 11, 0, 0, 0, 0, time.UTC), time.Date(2016, 7, 13, 23, 59, 0, 0, time.UTC)},
	// 	{4, "Comcast 2016", "Comcast, A good company. It's attended by the best and the brightest of the Go team and community.", "Denver", "http://comcast.com/", time.Date(2016, 7, 11, 0, 0, 0, 0, time.UTC), time.Date(2016, 7, 13, 23, 59, 0, 0, time.UTC)},
	// }
	OpenDB()
	// for i := 0; i < 4; i++ {
	// 	CreateEvent(&events[i])
	// }
	// ReadDB()
	ReadDBStream()
	defer db.Close()
}
