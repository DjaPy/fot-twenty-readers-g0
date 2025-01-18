package adapters

import (
	"context"
	"database/sql"
	"github.com/DjaPy/fot-twenty-readers-go/src/domain"
	"github.com/asdine/storm/v3"
	"github.com/gofrs/uuid/v5"
	"github.com/pkg/errors"
	"log"
	"time"
)

type PsalmReaderTGDB struct {
	Id         uuid.UUID `storm:"id"`
	Username   string    `storm:"index"`
	TelegramID int64     `storm:"index, unique"`
	Phone      string
	CreatedAt  time.Time `storm:"index"`
	UpdatedAt  time.Time
}

type PsalmReaderTGRepository struct {
	db *storm.DB
}

func NewPsalmReaderTGRepository(db *storm.DB) *PsalmReaderTGRepository {
	if db == nil {
		panic("missing db")
	}
	return &PsalmReaderTGRepository{db: db}
}

type sqlContextGetter interface {
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

func (pr PsalmReaderTGRepository) GetUserTG(ctx context.Context, id uuid.UUID) (*domain.PsalmReader, error) {
	return pr.getOrCreateHour(ctx, pr.db, id)
}

func (pr PsalmReaderTGRepository) getOrCreateHour(
	ctx context.Context,
	db *storm.DB,
	id uuid.UUID,
) (*domain.PsalmReader, error) {
	var dbPsalmReaderTG PsalmReaderTGDB
	err := db.One("ID", id, &dbPsalmReaderTG)
	if err != nil {
		log.Fatal(err)
	}

	if errors.Is(err, sql.ErrNoRows) {
		return nil, err
	} else if err != nil {
		return nil, errors.Wrap(err, "unable to get hour from db")
	}

	domainHour := domain.UnmarshallPsalmReader(
		dbPsalmReaderTG.Id,
		dbPsalmReaderTG.Username,
		dbPsalmReaderTG.TelegramID,
		dbPsalmReaderTG.Phone,
		dbPsalmReaderTG.CreatedAt,
		dbPsalmReaderTG.UpdatedAt,
	)

	return domainHour, nil
}
