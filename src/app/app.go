package app

import (
	"github.com/DjaPy/fot-twenty-readers-go/src/adapters"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	CreatePsalmReaderTG
	CreateCalendarOfReader
}

type Queries struct {
}
