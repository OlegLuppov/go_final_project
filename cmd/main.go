package main

import (
	"log"

	"github.com/OlegLuppov/go_final_project/config"
	"github.com/OlegLuppov/go_final_project/pkg/api"
	"github.com/OlegLuppov/go_final_project/pkg/db"
	"github.com/OlegLuppov/go_final_project/pkg/server"
)

func main() {
	env, err := config.LoadEnv()

	if err != nil {
		log.Fatalf("error LoadEnv: %s", err)
	}

	// Для прода не запкскать если нет секретного ключа для подписи токена
	// if len(env.SecretKey) == 0 {
	// 	log.Fatal("env SECRET_KEY is empty")
	// }

	schedulerDb, err := db.Connect(env.TodoDbFile)

	if err != nil {
		log.Fatalf("error db Connect: %s", err)
	}

	defer schedulerDb.Db.Close()

	router := api.RegisterHandlers(schedulerDb, env)

	err = server.Run(env, router)

	if err != nil {
		log.Fatalf("error server Run: %s", err)
	}

}
