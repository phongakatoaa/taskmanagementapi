package main

import (
	"database/sql"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"siransbach/taskmanagementapi/auth"
	"siransbach/taskmanagementapi/config"
	"siransbach/taskmanagementapi/handlers"
	"siransbach/taskmanagementapi/postgres"
)

func main() {
	setupLogger()

	pg, err := postgres.Connect(postgres.ConnectionString())
	if err != nil {
		log.Fatal().Err(err).Msg("error connecting to the postgres database")
	}
	defer func() {
		if err := pg.Close(); err != nil {
			log.Error().Err(err).Msg("error closing the postgres database")
		}
	}()

	app := setupFiberApp(pg)

	handlers.Setup(app, pg)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-shutdown
		log.Info().Msg("shutting down the server...")
		// graceful shutdown
		if err := app.Shutdown(); err != nil {
			log.Error().Err(err).Msg("error shutting down the server")
		}
	}()

	if err := app.Listen(config.HostPort()); err != nil {
		log.Fatal().Err(err).Msg("error starting server")
	}
}

func setupLogger() {
	logLevel, err := zerolog.ParseLevel(config.LogLevel())
	if err != nil {
		logLevel = zerolog.InfoLevel
	}
	log.Info().Msgf("setting log level to %s", logLevel)
	zerolog.SetGlobalLevel(logLevel)
	zerolog.TimeFieldFormat = time.RFC3339
}

func setupFiberApp(pg *sql.DB) *fiber.App {
	app := fiber.New(fiber.Config{})
	app.Use(logger.New(
		logger.Config{
			Next:         nil,
			Format:       "[${time}] ${status} - ${latency} ${method} ${path}\n",
			TimeFormat:   "15:04:05",
			TimeZone:     "Local",
			TimeInterval: 500 * time.Millisecond,
			Output:       os.Stderr,
		}),
	)
	app.Use(recover.New(recover.Config{EnableStackTrace: true}))
	app.Use(auth.NewMiddleware(pg))
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong")
	})

	return app
}
