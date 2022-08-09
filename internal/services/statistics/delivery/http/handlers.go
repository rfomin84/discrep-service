package statistics

import (
	"encoding/json"
	statistics "github.com/rfomin84/discrep-service/internal/services/statistics/domain"
	"io"
	"net/http"
)

func (d *Delivery) statistics(w http.ResponseWriter, r *http.Request) {

	type bodyParams struct {
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
		FeedIds   []int  `json:"feed_ids"`
	}

	var body bodyParams
	bytes, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Unmarshal(bytes, &body)

	result := d.statisticUseCase.GetStatistics(body.StartDate, body.EndDate, body.FeedIds)

	type response struct {
		Data []statistics.DetailedFeedStatistic `json:"data"`
	}

	resp := response{
		Data: result,
	}

	resultByte, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resultByte)
}
