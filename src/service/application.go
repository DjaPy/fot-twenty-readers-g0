package service

import (
	"github.com/DjaPy/fot-twenty-readers-go/src/adapters"
	"github.com/asdine/storm/v3"
	"github.com/asdine/storm/v3/codec/json"
	"log"
)

func NewApplication() {

	db, err := storm.Open("for-twenty-readers.db", storm.Codec(json.Codec))
	if err != nil {
		log.Fatalf("Could not open database: %v", err)
	}
	defer func(db *storm.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalf("Could not close database: %v", err)
		}
	}(db)

	psalmReaderTGRepository := adapters.NewPsalmReaderTGRepository(db)

}
