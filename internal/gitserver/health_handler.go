package gitserver

import "net/http"

type HealthResponse struct {
	Status string `json:"status"`
}

func healthHandler(w http.ResponseWriter, r *http.Request) {

	errJSON := HealthResponse{Status: "ok"}
	WriteIndentedJSON(w, errJSON, "", "  ")

}
