package server

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

const webDir = "./web"

func Run(port string) error {

	r := chi.NewRouter()

	fs := http.FileServer(http.Dir(webDir))
	r.Handle("/*", fs)

	return http.ListenAndServe(fmt.Sprintf(":%s", port), r)
}
