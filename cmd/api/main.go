package main

import (
	"log"
	"social/internal/env"
)

func main() {
	cfg := config{
		addr: env.GetString("ADDR", ":9091"),
	}
	app := &application{
		config: cfg,
		logger: log.Default(),
	}

	mux := app.mount()
	log.Fatal(app.run(mux))
}
