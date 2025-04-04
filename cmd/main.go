package main

import (
	"log"

	"go1f/config"
	"go1f/pkg/server"
)

func main() {
	env, err := config.LoadEnv()

	if err != nil {
		log.Fatalf("LoadEnv: %s", err)
	}

	err = server.Run(env.TodoPort)

	if err != nil {
		log.Fatalf("Error server Run: %s", err)
	}

}
