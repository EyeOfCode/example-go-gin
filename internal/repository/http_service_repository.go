package repository

import (
	"context"
	"net/http"
	"time"
)

type HttpServiceRepository  interface {
	Get(ctx context.Context, url string) error
}

type httpServiceRepository  struct {
	client *http.Client
}

func NewHttpServiceRepository () HttpServiceRepository  {
	return &httpServiceRepository {
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (s *httpServiceRepository ) Get(ctx context.Context, url string) error {
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
