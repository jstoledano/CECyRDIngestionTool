package main

import (
	"encoding/csv"
	"log"
	"os"

	"github.com/jstoledano/CECyRDIngestionTool/database"
)

func main() {
	db, err := database.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	err = database.InitTables(db)
	if err != nil {
		log.Fatal(err)
	}

	// Open the csv file
	f, err := os.Open("data/01.csv")
	if err != nil {
		log.Fatal(err)
	}
	csvReader := csv.NewReader(f)
	csvReader.Comma = '|'
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// var records []database.Record
	for _, row := range data {
		log.Println(row[0])
	}

}
