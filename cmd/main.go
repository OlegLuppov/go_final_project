package main

import (
	"fmt"
	"go1f/pkg/server"
	"log"
)

func main() {
	fmt.Println("Hello Final Project")
	err := server.Run()

	if err != nil {
		log.Fatalf("Error server Run: %s", err)
	}
}
