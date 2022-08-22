package feeds

type Feed struct {
	Id                 int      `json:"id"`
	Formats            []string `json:"placement_types"`
	UserId             int      `json:"userId"`
	ExternalStatistics bool     `json:"external_statistics"`
	RtbApiProviderId   int      `json:"rtb_api_provider_id"`
	IsDsp              bool     `json:"is_dsp"`
	TimezoneName       string   `json:"timezone_name"`
	TimezoneOffset     string   `json:"timezone_offset"`
}
