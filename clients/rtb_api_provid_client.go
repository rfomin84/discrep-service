package clients

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"net/http"
	"net/url"
	"time"
)

type RtbApiProvidClient struct {
	Client   *http.Client
	Host     string
	ApiToken string
}

func NewRtbApiProvidClient(cfg viper.Viper) *RtbApiProvidClient {
	httpClient := &http.Client{}
	transport := &http.Transport{}
	transport.MaxIdleConns = 20

	httpClient.Transport = transport
	//httpClient.Timeout = 5 * time.Second

	return &RtbApiProvidClient{
		Client:   httpClient,
		Host:     fmt.Sprintf("%s:%s", cfg.GetString("RTB_API_PROVID_HOST"), cfg.GetString("RTB_API_PROVID_PORT")),
		ApiToken: cfg.GetString("RTB_API_PROVID_TOKEN"),
	}
}

func (client *RtbApiProvidClient) GetStatistics(from, to time.Time, rtbApiProviderID string) (*http.Response, error) {
	ctx := context.Background()

	q := url.Values{}
	q.Set("from", from.Format("2006-01-02"))
	q.Set("to", to.Format("2006-01-02"))
	q.Set("rtb_api_provider_id", rtbApiProviderID)

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, client.Host+"/api/v1/stats?"+q.Encode(), nil)
	if err != nil {
		return nil, err
	}

	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", client.ApiToken))

	return client.Client.Do(request)
}
