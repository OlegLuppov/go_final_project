package server

import (
	"fmt"
	"net/http"

	"github.com/OlegLuppov/go_final_project/config"
	"github.com/go-chi/chi"
)

const webDir = "./web"

func Run(env config.Environment, router *chi.Mux) error {

	fs := http.FileServer(http.Dir(webDir))
	router.Handle("/*", fs)

	return http.ListenAndServe(fmt.Sprintf(":%s", env.TodoPort), router)
}
