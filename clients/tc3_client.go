package clients

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"net/http"
)

type Tc3Client struct {
	Client   *http.Client
	Host     string
	ApiToken string
}

func New(cfg *viper.Viper) *Tc3Client {
	httpClient := &http.Client{}
	transport := &http.Transport{}
	transport.MaxIdleConns = 20

	httpClient.Transport = transport

	return &Tc3Client{
		Client:   httpClient,
		Host:     fmt.Sprintf("%s:%s", cfg.GetString("TC3_HOST"), cfg.GetString("TC3_PORT")),
		ApiToken: cfg.GetString("TC3_API_TOKEN"),
	}
}

func (httpClient Tc3Client) GetFeeds() (*http.Response, error) {
	//ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	//defer cancel()

	ctx := context.Background()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, httpClient.Host+"/api/internal/v1/feeds-discrep-list", nil)
	if err != nil {
		return nil, err
	}

	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", httpClient.ApiToken))

	return httpClient.Client.Do(request)
}
