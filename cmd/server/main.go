package main

import (
	"context"
	"os"
	"slash10k/internal/db"
	"slash10k/internal/handler"
	"slash10k/internal/utils"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/time/rate"
)

func main() {
	setupLogger()

	err := db.Migrate(context.Background(), utils.DefaultConfig().ConnectionString)
	if err != nil {
		log.Fatal().Err(err).Msg("could not run database migrations")
		return
	}

	d := db.NewDatabase()

	e := echo.New()

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Recover())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(
		rate.Limit(20),
	)))

	log.Debug().Msg("setting up routes")

	e.GET("/", handler.ServeFrontend)
	e.Static("/", "assets")

	api := e.Group("/api")
	api.GET("/debt", handler.AllDebts(d))
	api.GET("/debt/:player", handler.GetDebt)
	api.POST("/debt/:player/:amount", handler.AddDebt(d))

	api.GET("/journal/:player", handler.GetJournalEntries(d))

	admin := api.Group("/admin")

	if adminPw == "disabled" {
		log.Warn().Msg("admin password is disabled")
	} else {
		log.Debug().Msg("admin password is enabled")
		pw := os.Getenv("ADMIN_PASSWORD")
		if pw == "" {
			log.Fatal().Msg("ADMIN_PASSWORD not set")
			return
		}
		admin.Use(middleware.BasicAuth(func(user, pass string, context echo.Context) (bool, error) {
			return user == "admin" && pass == pw, nil
		}))
	}

	admin.POST("/player/:name", handler.AddPlayer(d))
	admin.DELETE("/player/:name", handler.DeletePlayer(d))
	admin.POST("/char", handler.AddChar)
	admin.DELETE("/char/:name", handler.DeleteChar)

	log.Info().Msgf("server started on port %v", 3000)
	e.Logger.Fatal(e.Start(":3000"))
}

func setupLogger() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel != "" {
		l, err := zerolog.ParseLevel(strings.ToLower(logLevel))
		if err != nil {
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
		} else {
			zerolog.SetGlobalLevel(l)
		}
	} else {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	output := zerolog.ConsoleWriter{Out: os.Stderr}
	log.Logger = log.Output(output)
}

// go build -ldflags "-X main.adminPwDisabled=disabled" -o slash10k
var adminPw string
