package model

import (
	"context"
	"github.com/google/uuid"
	"time"
)

type Media struct {
	Filename     string    `json:"filename"`
	Date         string    `json:"date"`
	Url          string    `json:"url"`
	LastModified time.Time `json:"last_modified"`
}

type MediaUsecase interface {
	Download(ctx context.Context, userID uuid.UUID, date time.Time, filename string) ([]byte, error)
	GetAPOD(ctx context.Context) (*Media, error)
	List(ctx context.Context, userID uuid.UUID) ([]Media, error)
	ListByDate(ctx context.Context, userID uuid.UUID, date time.Time) ([]Media, error)
	Store(ctx context.Context, userID uuid.UUID, date time.Time, filename string, bytes []byte) error
}
