package adapters

import (
	"errors"
	"github.com/DjaPy/fot-twenty-readers-go/src/domain"
	"github.com/asdine/storm/v3"
	"github.com/gofrs/uuid/v5"
	"log"
	"time"
)

type CalendarOfReaderDB struct {
	ID        uuid.UUID `storm:"id"`
	Calendar  domain.CalendarMap
	CreatedAt time.Time `storm:"index"`
	UpdatedAt time.Time
}

type CalendarOfReaderRepository struct {
	db *storm.DB
}

func NewCalendarOfReaderRepository(db *storm.DB) *CalendarOfReaderRepository {
	if db == nil {
		log.Fatal("missing db")
	}
	return &CalendarOfReaderRepository{db: db}
}

func (cr CalendarOfReaderRepository) GetCalendar(id uuid.UUID) (*CalendarOfReaderDB, error) {
	var calendarOfReader CalendarOfReaderDB
	err := cr.db.One("ID", id, &calendarOfReader)
	if err != nil {
		return nil, err
	}
	return &calendarOfReader, nil
}

func (cr CalendarOfReaderRepository) CreateCalendarOfReader(
	calendarOfReader *domain.CalendarOfReader,
) error {
	err := cr.db.Save(&calendarOfReader)
	if err != nil {
		if errors.Is(err, storm.ErrAlreadyExists) {
			return err
		}
	}
	return nil
}
