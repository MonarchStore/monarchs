package server

import (
	"encoding/json"
	"net/http"
)

type Metric struct {
	Name string
	Time string
	Values map[string]string
}

func (s *httpServer) healthCheck(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (s *httpServer) doMetrics(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	m := Metric{"Foobar", "11:39", map[string]string{"key": "value"}}
	js, err := json.Marshal(m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
