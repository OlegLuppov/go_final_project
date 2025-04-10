package server

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

const webDir = "./web"

func Run(port string, router *chi.Mux) error {

	fs := http.FileServer(http.Dir(webDir))
	router.Handle("/*", fs)

	return http.ListenAndServe(fmt.Sprintf(":%s", port), router)
}
