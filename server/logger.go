package server

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

func Logger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
	})
}
