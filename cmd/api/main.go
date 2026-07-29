package main

import "log"

func main() {
	cfg := config{
		addr: ":9090",
	}
	app := &application{
		config: cfg,
	}

	log.Fatal(app.run())
}
