// @title           Social API
// @version         1.0.0
// @description     REST API for the social backend built with Go.
//
// @host            localhost:9090
// @BasePath        /v1
//
// @contact.name    Fernando Piedade

package main

import (
	"log"
	_ "social/docs"
	"social/internal/db"
	"social/internal/env"
	"social/internal/store"
)

func main() {
	cfg := config{
		addr: env.GetString("ADDR", ":9090"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://postgres:postgres@localhost/social?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
	}

	db, err := db.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)

	if err != nil {
		log.Panic(err)
	}

	defer db.Close()
	log.Println("Database Connection Pool Established")

	store := store.NewStorage(db)

	app := &application{
		config: cfg,
		logger: log.Default(),
		store:  store,
	}

	mux := app.mount()
	log.Fatal(app.run(mux))
}
