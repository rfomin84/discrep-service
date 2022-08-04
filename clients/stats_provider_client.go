package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"net/http"
	"time"
)

type StatsProviderClient struct {
	Client   *http.Client
	Host     string
	ApiToken string
}

type data struct {
	StartDate   time.Time `json:"startDate"`
	EndDate     time.Time `json:"endDate"`
	Feeds       []int     `json:"feeds"`
	BillingType string    `json:"billingType"`
	Timeframe   string    `json:"timeframe"`
}

func (d data) MarshalJSON() ([]byte, error) {
	type dataAlias data

	aliasValue := struct {
		dataAlias
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}{
		dataAlias: dataAlias(d),
		StartDate: d.StartDate.Format("2006-01-02 15:04:05"),
		EndDate:   d.EndDate.Format("2006-01-02 15:04:05"),
	}

	return json.Marshal(aliasValue)
}

func (d *data) UnmarshalJSON(dataBytes []byte) error {
	type dataAlias data

	aliasValue := &struct {
		*dataAlias
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}{
		dataAlias: (*dataAlias)(d),
	}

	if err := json.Unmarshal(dataBytes, aliasValue); err != nil {
		return err
	}

	d.StartDate, _ = time.Parse("2006-01-02 15:04:05", aliasValue.StartDate)
	d.EndDate, _ = time.Parse("2006-01-02 15:04:05", aliasValue.EndDate)

	return nil
}

func NewStatsProviederClient(cfg *viper.Viper) *StatsProviderClient {
	httpClient := &http.Client{}
	transport := &http.Transport{}
	transport.MaxIdleConns = 20

	httpClient.Transport = transport

	return &StatsProviderClient{
		Client:   httpClient,
		Host:     fmt.Sprintf("%s:%s", cfg.GetString("STATS_PROVIDER_HOST"), cfg.GetString("STATS_PROVIDER_PORT")),
		ApiToken: cfg.GetString("STATS_PROVIDER_API_TOKEN"),
	}
}

func (spc *StatsProviderClient) GetStatistics(startDate, endDate time.Time, billingType, timeframe string, feedIds []int) (*http.Response, error) {
	ctx := context.Background()

	dataBody := data{
		StartDate:   startDate,
		EndDate:     endDate,
		Feeds:       feedIds,
		BillingType: billingType,
		Timeframe:   timeframe,
	}

	dataByte, err := json.Marshal(dataBody)
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, spc.Host+"/api/v1/billing-stats-by-feeds", bytes.NewBuffer(dataByte))

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", spc.ApiToken))

	return spc.Client.Do(request)
}
