package domain

import (
	"github.com/gofrs/uuid/v5"
	"time"
)

type PsalmReader struct {
	ID         uuid.UUID
	Username   string
	TelegramID int64
	Phone      string
	CalendarID uuid.UUID
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func NewPsalmReader(username string, telegramID int64, phone string) (*PsalmReader, error) {
	ID, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	createdAt := time.Now()
	updatedAt := time.Now()
	return &PsalmReader{
		ID:         ID,
		Username:   username,
		TelegramID: telegramID,
		Phone:      phone,
		CalendarID: uuid.Nil,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}, nil
}

func UnmarshallPsalmReader(
	id uuid.UUID,
	username string,
	telegramID int64,
	phone string,
	CalendarID uuid.UUID,
	createdAt time.Time,
	updatedAt time.Time,
) *PsalmReader {
	return &PsalmReader{
		ID:         id,
		Username:   username,
		TelegramID: telegramID,
		Phone:      phone,
		CalendarID: CalendarID,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}
}
