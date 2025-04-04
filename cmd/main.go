package main

import (
	"fmt"
	"log"

	"github.com/OlegLuppov/go_final_project/pkg/server"
)

func main() {
	fmt.Println("Hello Final Project")
	err := server.Run()

	if err != nil {
		log.Fatalf("Error server Run: %s", err)
	}
}
