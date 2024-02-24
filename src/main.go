package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/time/rate"
	"os"
	"scurvy10k/src/handler"
)

func main() {
	setupLogger()

	e := echo.New()

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Recover())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(
		rate.Limit(20),
	)))

	log.Debug().Msg("setting up routes")

	e.GET("/", handler.ServeFrontend)

	api := e.Group("/api")
	api.GET("/debt/:player", handler.GetDebt)
	api.POST("/debt/:player/:amount", handler.AddDebt)

	admin := api.Group("/admin")
	admin.POST("/player/:name", handler.AddPlayer)
	admin.DELETE("/player/:name", handler.DeletePlayer)
	admin.POST("/char", handler.AddChar)
	admin.DELETE("/char/:name", handler.DeleteChar)

	log.Info().Msgf("server started on port %v", 3000)
	e.Logger.Fatal(e.Start(":3000"))
}

func setupLogger() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	output := zerolog.ConsoleWriter{Out: os.Stderr}
	log.Logger = log.Output(output)
}
