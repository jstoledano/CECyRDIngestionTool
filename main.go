package main

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {
	db, err := gorm.Open(sqlite.Open("cecyrd.db"), &gorm.Config{})
	log.Println(db, err)
}
