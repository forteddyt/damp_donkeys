package main

import (
	"log"
	
	"github.com/damp_donkeys/internal/app/router"
)

func main() {
	log.Print("Starting server...")

	router.Setup()
}

