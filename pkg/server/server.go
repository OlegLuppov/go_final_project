package server

import (
	"fmt"
	"log"
	"net/http"

	"go1f/config"

	"github.com/go-chi/chi"
)

const webDir = "./web"

func Run() error {
	env, err := config.LoadEnv()

	if err != nil {
		log.Fatalf("LoadEnv: %s", err)
	}

	r := chi.NewRouter()

	fs := http.FileServer(http.Dir(webDir))
	r.Handle("/*", fs)

	return http.ListenAndServe(fmt.Sprintf(":%s", env.TodoPort), r)
}
