package api

import (
	"net/http"

	"go1f/pkg/dateutil"

	"github.com/go-chi/chi"
)

func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	queryNow := r.FormValue("now")
	queryDate := r.FormValue(("date"))
	queryRepeat := r.FormValue(("repeat"))

	nextDate, err := dateutil.NextDate(queryNow, queryDate, queryRepeat)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(nextDate))
}

func RegisterHandlers() *chi.Mux {
	router := chi.NewMux()

	router.Get("/api/nextdate", NextDateHandler)

	return router
}
