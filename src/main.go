package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/time/rate"
	"io"
	"os"
	"scurvy10k/src/handler"
	"text/template"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	setupLogger()

	e := echo.New()

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Recover())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(
		rate.Limit(20),
	)))

	log.Debug().Msg("setting up routes")

	t := &Template{
		templates: template.Must(template.ParseGlob("templ/*.html")),
	}
	e.Renderer = t

	e.GET("/", handler.ServeFrontend)

	api := e.Group("/api")
	api.GET("/debt", handler.Debt)

	log.Info().Msgf("server started on port %v", 3000)
	e.Logger.Fatal(e.Start(":3000"))
}

func setupLogger() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}
