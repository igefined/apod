package usecase

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/igilgyrg/betera-test/internal/model"
	"github.com/igilgyrg/betera-test/pkg/nasa"
	"github.com/igilgyrg/betera-test/pkg/storage/s3"
	"time"
)

type mediaUsecase struct {
	nasaClient nasa.NASAClient
	storage    s3.S3Storage
	ctxTimeout time.Duration
}

func NewMediaUsecase(nasaClient nasa.NASAClient, s3Storage s3.S3Storage, ctxTimeout time.Duration) model.MediaUsecase {
	return &mediaUsecase{nasaClient: nasaClient, storage: s3Storage, ctxTimeout: ctxTimeout}
}

func (m mediaUsecase) GetAPOD(c context.Context) (*model.Media, error) {
	_, cancel := context.WithTimeout(c, m.ctxTimeout)
	defer cancel()

	mediaAPOD, err := m.nasaClient.APOD(time.Now().UTC())
	if err != nil {
		return nil, err
	}

	lastModified, err := time.Parse("2006-01-02", mediaAPOD.Date)
	if err != nil {
		return nil, err
	}

	return &model.Media{
		Filename:     mediaAPOD.Title,
		Url:          mediaAPOD.HDUrl,
		Date:         time.Now().Format("2006-01-02"),
		LastModified: lastModified,
	}, nil
}

func (m mediaUsecase) Download(c context.Context, userID uuid.UUID, date time.Time, filename string) ([]byte, error) {
	_, cancel := context.WithTimeout(c, m.ctxTimeout)
	defer cancel()

	path := fmt.Sprintf("%s/%s/%s", userID.String(), date.Format("2006-01-02"), filename)
	bytes, err := m.storage.Download(path)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func (m mediaUsecase) List(c context.Context, userID uuid.UUID) ([]model.Media, error) {
	_, cancel := context.WithTimeout(c, m.ctxTimeout)
	defer cancel()

	medias, err := m.storage.List(userID.String())
	if err != nil {
		return nil, err
	}

	result := make([]model.Media, len(medias))
	for i := range medias {
		m := medias[i]
		result[i] = model.Media{
			Filename:     m.Filename,
			Url:          m.Url,
			Date:         m.Date,
			LastModified: m.LastModified,
		}
	}

	return result, nil
}

func (m mediaUsecase) ListByDate(c context.Context, userID uuid.UUID, date time.Time) ([]model.Media, error) {
	_, cancel := context.WithTimeout(c, m.ctxTimeout)
	defer cancel()

	filename := fmt.Sprintf("%s/%s", userID.String(), date.Format("2006-01-02"))
	medias, err := m.storage.List(filename)
	if err != nil {
		return nil, err
	}

	result := make([]model.Media, len(medias))
	for i := range medias {
		m := medias[i]
		result[i] = model.Media{
			Filename:     m.Filename,
			Url:          m.Url,
			Date:         m.Date,
			LastModified: m.LastModified,
		}
	}

	return result, nil
}

func (m mediaUsecase) Store(c context.Context, userID uuid.UUID, date time.Time, filename string, bytes []byte) error {
	_, cancel := context.WithTimeout(c, m.ctxTimeout)
	defer cancel()

	return m.storage.Store(fmt.Sprintf("%s/%s/%s", userID.String(), date.Format("2006-01-02"), filename), bytes)
}
