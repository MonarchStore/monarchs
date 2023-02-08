package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
)

type Metrics struct {
	Store  map[string]string
	Memory map[string]string
}

func (s *httpServer) collectMetrics() *Metrics {

	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	mm := make(map[string]string)
	mm["MemAlloc"] = fmt.Sprintf("%v", mem.Alloc)
	mm["MemTotalAlloc"] = fmt.Sprintf("%v", mem.TotalAlloc)
	mm["MemSys"] = fmt.Sprintf("%v", mem.Sys)
	mm["MemNumGC"] = fmt.Sprintf("%v", mem.NumGC)

	ss := make(map[string]string)
	ss["Size"] = fmt.Sprintf("%v", len(s.storeMap))

	return &Metrics{
		Store:  ss,
		Memory: mm,
	}
}

func (s *httpServer) healthCheck(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (s *httpServer) doMetrics(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)

	m := s.collectMetrics()
	js, err := json.Marshal(*m)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
