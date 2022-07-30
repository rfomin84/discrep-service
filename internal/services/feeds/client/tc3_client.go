package client

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"net/http"
	"time"
)

type HttpClient struct {
	Client   *http.Client
	Host     string
	ApiToken string
}

func New(cfg *viper.Viper) *HttpClient {
	httpClient := &http.Client{}
	transport := &http.Transport{}
	transport.MaxIdleConns = 20

	httpClient.Transport = transport

	return &HttpClient{
		Client:   httpClient,
		Host:     fmt.Sprintf("%s:%s", cfg.GetString("TC3_HOST"), cfg.GetString("TC3_PORT")),
		ApiToken: cfg.GetString("TC3_API_TOKEN"),
	}
}

func (httpClient HttpClient) GetFeeds() (*http.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, httpClient.Host+"/feeds", nil)
	if err != nil {
		return nil, err
	}

	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", httpClient.ApiToken))

	return httpClient.Client.Do(request)
}
