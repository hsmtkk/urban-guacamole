package main

import (
	"fmt"
	"html/template"
	"io"
	"log"

	"github.com/gorilla/sessions"
	"github.com/hsmtkk/urban-guacamole/env"
	"github.com/hsmtkk/urban-guacamole/handler"
	storearticle "github.com/hsmtkk/urban-guacamole/store/article"
	storesession "github.com/hsmtkk/urban-guacamole/store/session"
	"github.com/labstack/echo-contrib/session"
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

	port, err := env.GetPort()
	if err != nil {
		sugar.Fatal(err)
	}

	articleStore := storearticle.MemoryImpl()
	sessionStore := storesession.MemoryImpl()

	hdl := handler.New(sugar, articleStore, sessionStore)

	e := echo.New()

	t := &Template{
		templates: template.Must(template.ParseGlob("template/*.html")),
	}
	e.Renderer = t

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))

	e.GET("/login", hdl.GetLogin)
	e.POST("/login", hdl.PostLogin)

	loginRequiredGroup := e.Group("/auth", hdl.LoginRequired)

	loginRequiredGroup.GET("/logout", hdl.GetLogout)
	loginRequiredGroup.GET("/articles", hdl.GetArticles)
	loginRequiredGroup.GET("/articles/create", hdl.GetArticlesCreate)
	loginRequiredGroup.POST("/articles/create", hdl.PostArticlesCreate)
	loginRequiredGroup.GET("/articles/:id/read", hdl.GetArticlesRead)
	loginRequiredGroup.GET("/articles/:id/update", hdl.GetArticlesUpdate)
	loginRequiredGroup.POST("/articles/:id/update", hdl.PostArticlesUpdate)
	loginRequiredGroup.POST("/articles/:id/delete", hdl.PostArticlesDelete)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port)))
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
