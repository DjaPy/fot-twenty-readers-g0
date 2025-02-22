package domain

import (
	"github.com/gofrs/uuid/v5"
	"time"
)

type CalendarMap map[string]map[string]string

type CalendarOfReader struct {
	ID        uuid.UUID `storm:"id"`
	Calendar  CalendarMap
	CreatedAt time.Time `storm:"index"`
	UpdatedAt time.Time
}
