package main

import (
	"fmt"
	"io"
	"log"
	"text/template"

	"github.com/hsmtkk/urban-guacamole/controller"
	"github.com/hsmtkk/urban-guacamole/entryrepo"
	"github.com/hsmtkk/urban-guacamole/env"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()
	sugar := logger.Sugar()

	entryRepo, err := entryrepo.NewFileImpl(sugar, "entryrepo/entryrepo.json")
	if err != nil {
		sugar.Fatal(err)
	}

	port, err := env.GetPort()
	if err != nil {
		sugar.Fatal(err)
	}

	ctrl := controller.New(sugar, entryRepo)

	// Echo instance
	e := echo.New()
	t := &Template{
		templates: template.Must(template.ParseGlob("template/*.html")),
	}
	e.Renderer = t

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", ctrl.ShowEntries)
	e.POST("/entries", ctrl.AddEntry)
	e.GET("/entries/new", ctrl.NewEntry)
	e.GET("/entries/:id", ctrl.ShowEntry)
	e.GET("/entries/:id/edit", ctrl.EditEntry)
	e.POST("/entries/:id/update", ctrl.UpdateEntry)
	e.POST("/entries/:id/delete", ctrl.DeleteEntry)

	// Start server
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port)))
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
