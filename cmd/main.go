package main

import (
	"log"

	"go1f/config"
	"go1f/pkg/server"

	"github.com/OlegLuppov/go_final_project/pkg/db"
)

func main() {
	env, err := config.LoadEnv()

	if err != nil {
		log.Fatalf("LoadEnv: %s", err)
	}

	schedulerDb, err := db.Connect(env.TodoDbFile)

	if err != nil {
		log.Fatalf("error db Connect: %s", err)
	}

	defer schedulerDb.Close()

	err = server.Run(env.TodoPort)

	if err != nil {
		log.Fatalf("Error server Run: %s", err)
	}

}
