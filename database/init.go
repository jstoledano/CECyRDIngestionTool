package database

import (
	"log"
	"os"
	"sync"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB
var once sync.Once
var err error

// Function InitDB() initializes the database connection
// using a singleton pattern y returns a gorm.DB object
func InitDB() (*gorm.DB, error) {

	// Using sync.Do we make sure that the connection to the database
	// is initialized only once, avoiding multiple connections
	once.Do(func() {
		if _, err := os.Stat("cecyrd.db"); os.IsNotExist(err) {
			db, err = gorm.Open(sqlite.Open("cecyrd.db"), &gorm.Config{})
			if err != nil {
				log.Println("failed to connect database", err)
				return
			}
			log.Println("Database created")
		} else {
			db, err = gorm.Open(sqlite.Open("cecyrd.db"), &gorm.Config{})
			if err != nil {
				log.Println("failed to connect database", err)
				return
			}
			log.Println("Database connected")
		}
	})

	return db, nil
}

// Function InitTables(db *gorm.DB) check if the Record table exists
// and if not, it creates it.
func InitTables(db *gorm.DB) error {
	if err := db.AutoMigrate(&Record{}); err != nil {
		return err
	}
	return nil
}
