package main

import (
	"os"
	"scurvy10k/src/handler"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/time/rate"
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
	e.Static("/", "assets")

	api := e.Group("/api")
	api.GET("/debt", handler.AllDebts)
	api.GET("/debt/:player", handler.GetDebt)
	api.POST("/debt/:player/:amount", handler.AddDebt)

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

	admin.POST("/player/:name", handler.AddPlayer)
	admin.DELETE("/player/:name", handler.DeletePlayer)
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

// go build -ldflags "-X main.adminPwDisabled=disabled" -o scurvy10k
var adminPw string
