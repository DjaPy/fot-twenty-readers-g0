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
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func UnmarshallPsalmReader(
	id uuid.UUID,
	username string,
	telegramID int64,
	phone string,
	createdAt time.Time,
	updatedAt time.Time,
) *PsalmReader {
	return &PsalmReader{
		ID:         id,
		Username:   username,
		TelegramID: telegramID,
		Phone:      phone,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}
}
