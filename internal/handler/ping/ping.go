package ping

import (
	"encoding/json"
	"net/http"
)

func Ping(w http.ResponseWriter, _ *http.Request) {
	type data struct {
		RequestParam string `json:"request_param"`
		Status       string `json:"status"`
		ErrorMessage string `json:"error_message"`
		Data         string `json:"data"`
		Next         string `json:"next"`
		Version      struct {
			Code string `json:"code"`
			Name string `json:"name"`
		}
	}

	res := data{
		Status: "success",
		Data:   "health",
		Version: struct {
			Code string `json:"code"`
			Name string `json:"name"`
		}{
			Code: "1",
			Name: "0.1.0",
		},
	}

	bres, _ := json.Marshal(res)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bres)
}
