package server

import (
	"fmt"
	"net/http"

	"go1f/config"

	"github.com/go-chi/chi"
)

const webDir = "./web"

func Run() error {
	env := config.LoadEnv()

	r := chi.NewRouter()

	fs := http.FileServer(http.Dir(webDir))
	r.Handle("/*", fs)

	fmt.Println(env.TodoPort)
	return http.ListenAndServe(fmt.Sprintf(":%s", env.TodoPort), r)
}
