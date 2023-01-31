package handlers

import (
	"github.com/tiraill/go_collect_metrics/internal/storage"
	"net/http"
)

func GetPingHandler(db storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "text/plain")
		ok := db.Ping()
		if ok {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("true"))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("false"))
		}
	}
}
