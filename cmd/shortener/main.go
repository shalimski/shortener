package main

import (
	"log"

	"github.com/shalimski/shortener/internal/app"

	"github.com/shalimski/shortener/config"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	app.Run(cfg)
}
