package service

import (
	"context"
	"net/http"
	"time"
)

type HttpService  interface {
	Get(ctx context.Context, url string) error
}

type httpService  struct {
	client *http.Client
}

func NewHttpService () HttpService  {
	return &httpService {
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (s *httpService ) Get(ctx context.Context, url string) error {
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
