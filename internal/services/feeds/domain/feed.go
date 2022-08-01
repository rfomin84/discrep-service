package feeds

type Feed struct {
	Id                 int      `json:"id"`
	Formats            []string `json:"formats"`
	UserId             int      `json:"userId"`
	ExternalStatistics bool     `json:"external_statistics"`
	IsDsp              bool     `json:"is_dsp"`
	TimezoneName       string   `json:"timezone_name"`
	TimezoneOffset     string   `json:"timezone_offset"`
}
