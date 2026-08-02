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
	"context"
	"io"
	"log/slog"
	"os"
	_ "social/docs"
	"social/internal/db"
	"social/internal/env"
	"social/internal/store"

	"github.com/natefinch/lumberjack"
)

func main() {
	// criar pasta logs se não existir
	if err := os.MkdirAll("logs", 0755); err != nil {
		panic(err)
	}

	logRotator := &lumberjack.Logger{
		Filename:   "logs/social-app.log",
		MaxSize:    1000, // roda ao atingir 1000 MB
		MaxBackups: 30,   // mantém no máximo 30 ficheiros antigos
		MaxAge:     30,   // apaga ficheiros com mais de 30 dias
		Compress:   true, // comprime os ficheiros antigos em .gz
	}
	defer logRotator.Close()

	multiWriter := io.MultiWriter(os.Stdout, logRotator)

	logger := slog.New(slog.NewJSONHandler(multiWriter, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	// cfg := config{
	// 	addr: env.GetString("ADDR", ":9090"),
	// 	db: dbConfig{
	// 		addr:         env.GetString("DB_ADDR", "postgres://postgres:postgres@localhost/social?sslmode=disable"),
	// 		maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
	// 		maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
	// 		maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
	// 	},
	// }

	cfg := config{
		addr:         env.GetString("ADDR", ":9090"),
		metricsAddr:  env.GetString("METRICS_ADDR", "127.0.0.1:9091"),
		otelEndpoint: env.GetString("OTEL_ENDPOINT", "localhost:4318"),
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
		logger.Error("database connection failed", "error", err)
		os.Exit(1)
	}
	defer db.Close()
	logger.Info("database connection pool established")

	// inicializar OTel
	ctx := context.Background()
	tp, err := newTracerProvider(ctx, cfg.otelEndpoint)
	if err != nil {
		logger.Error("failed to initialize tracer", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			logger.Error("tracer shutdown error", "error", err)
		}
	}()
	logger.Info("tracer initialized", "endpoint", cfg.otelEndpoint)

	store := store.NewStorage(db)

	app := &application{
		config:  cfg,
		logger:  logger,
		store:   store,
		metrics: newMetrics(),
	}

	mux := app.mount()
	if err := app.run(mux); err != nil {
		logger.Error("server error", "error", err)
		os.Exit(1)
	}
}
