package balance_history

import (
	"encoding/json"
	"net/http"
)

func (d *BalanceHistoryDelivery) reserveBalance(w http.ResponseWriter, r *http.Request) {
	result, err := d.balanceHistoryUseCase.ReservedFeedBalance()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resultByte, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resultByte)
}
