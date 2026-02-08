package main

import (
	"log"

	"github.com/curtisnewbie/nota/internal/app"
)

func main() {
	notaApp, err := app.NewApp()
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}
	
	notaApp.Run()
}