package clients

import "net/http"

type StatsProviderClient struct {
	Client   *http.Client
	Host     string
	ApiToken string
}

func (spc *StatsProviderClient) GetStatistics(format string, feedIds []int) interface{} {
	panic("реализуй меня")
}
