package repository

import (
	"context"
	"net/http"
	"time"
)

type ServiceRepository interface {
	Get(ctx context.Context, url string) error
}

type serviceRepository struct {
	client *http.Client
}

func NewServiceRepository() ServiceRepository {
	return &serviceRepository{
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (s *serviceRepository) Get(ctx context.Context, url string) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	switch res.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusNotFound:
		return nil
	default:
		return err
	}
}