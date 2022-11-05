package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/hsmtkk/urban-guacamole/record"
	storearticle "github.com/hsmtkk/urban-guacamole/store/article"
	storesession "github.com/hsmtkk/urban-guacamole/store/session"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

const (
	constSessionName = "session"
	constSessionID   = "hogefuga"
)

type MyHandler interface {
	GetLogin(ectx echo.Context) error
	PostLogin(ectx echo.Context) error
	GetLogout(ectx echo.Context) error
	GetArticles(ectx echo.Context) error
	GetArticlesCreate(ectx echo.Context) error
	PostArticlesCreate(ectx echo.Context) error
	GetArticlesRead(ectx echo.Context) error
	GetArticlesUpdate(ectx echo.Context) error
	PostArticlesUpdate(ectx echo.Context) error
	PostArticlesDelete(ectx echo.Context) error
	LoginRequired(next echo.HandlerFunc) echo.HandlerFunc
}

type handlerImpl struct {
	sugar        *zap.SugaredLogger
	articleStore storearticle.ArticleStore
	sessionStore storesession.SessionStore
}

func New(sugar *zap.SugaredLogger, articleStore storearticle.ArticleStore, sessionStore storesession.SessionStore) MyHandler {
	return &handlerImpl{sugar, articleStore, sessionStore}
}

func (h *handlerImpl) GetLogin(ectx echo.Context) error {
	h.sugar.Info("GetLogin")
	return ectx.Render(http.StatusOK, "login", nil)
}

func (h *handlerImpl) PostLogin(ectx echo.Context) error {
	h.sugar.Info("PostLogin")
	userID := ectx.FormValue("userid")
	sessionObj := record.NewSession(userID)
	if err := h.sessionStore.Set(ectx.Request().Context(), sessionObj); err != nil {
		return err
	}
	sess, err := session.Get(constSessionName, ectx)
	if err != nil {
		return fmt.Errorf("session cookie does not exist")
	}
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
	}
	sess.Values[constSessionID] = sessionObj.SessionID
	if err := sess.Save(ectx.Request(), ectx.Response()); err != nil {
		return err
	}
	return ectx.Redirect(http.StatusFound, "/auth/articles")
}

func (h *handlerImpl) GetLogout(ectx echo.Context) error {
	h.sugar.Info("GetLogout")
	sessionID, _, err := h.parseSession(ectx)
	if err != nil {
		return err
	}
	if err := h.sessionStore.Delete(ectx.Request().Context(), sessionID); err != nil {
		return err
	}
	return ectx.Redirect(http.StatusFound, "/login")
}

func (h *handlerImpl) GetArticles(ectx echo.Context) error {
	h.sugar.Info("GetArticles")
	_, userID, err := h.parseSession(ectx)
	if err != nil {
		return err
	}
	articles, err := h.articleStore.GetByUserID(ectx.Request().Context(), userID)
	if err != nil {
		return err
	}
	return ectx.Render(http.StatusOK, "articles", articles)
}

func (h *handlerImpl) GetArticlesCreate(ectx echo.Context) error {
	h.sugar.Info("GetArticlesCreate")
	return ectx.Render(http.StatusOK, "articlescreate", nil)
}

func (h *handlerImpl) PostArticlesCreate(ectx echo.Context) error {
	h.sugar.Info("PostArticlesCreate")
	_, userID, err := h.parseSession(ectx)
	if err != nil {
		return err
	}
	title := ectx.FormValue("title")
	text := ectx.FormValue("text")
	article := record.NewArticle(userID, title, text)
	if err := h.articleStore.Set(ectx.Request().Context(), article); err != nil {
		return err
	}
	return ectx.Redirect(http.StatusFound, "/auth/articles")
}

func (h *handlerImpl) GetArticlesRead(ectx echo.Context) error {
	h.sugar.Info("GetArticlesRead")
	articleID := ectx.Param("id")
	_, userID, err := h.parseSession(ectx)
	if err != nil {
		return err
	}
	article, err := h.articleStore.Get(ectx.Request().Context(), articleID)
	if err != nil {
		return err
	}
	if userID != article.UserID {
		return ectx.String(http.StatusForbidden, "you are not the owner of this article")
	}
	return ectx.Render(http.StatusOK, "articlesread", article)
}

func (h *handlerImpl) GetArticlesUpdate(ectx echo.Context) error {
	h.sugar.Info("GetArticlesUpdate")
	articleID := ectx.Param("id")
	_, userID, err := h.parseSession(ectx)
	if err != nil {
		return err
	}
	article, err := h.articleStore.Get(ectx.Request().Context(), articleID)
	if err != nil {
		return err
	}
	if userID != article.UserID {
		return ectx.String(http.StatusForbidden, "you are not the owner of this article")
	}
	return ectx.Render(http.StatusOK, "articlesupdate", article)
}

func (h *handlerImpl) PostArticlesUpdate(ectx echo.Context) error {
	h.sugar.Info("PostArticlesUpdate")
	articleID := ectx.Param("id")
	_, userID, err := h.parseSession(ectx)
	if err != nil {
		return err
	}
	article, err := h.articleStore.Get(ectx.Request().Context(), articleID)
	if err != nil {
		return err
	}
	if userID != article.UserID {
		return ectx.String(http.StatusForbidden, "you are not the owner of this article")
	}
	newTitle := ectx.FormValue("title")
	newText := ectx.FormValue("text")
	article.Title = newTitle
	article.Text = newText
	if err := h.articleStore.Set(ectx.Request().Context(), article); err != nil {
		return err
	}
	return ectx.Redirect(http.StatusFound, "/auth/articles")
}

func (h *handlerImpl) PostArticlesDelete(ectx echo.Context) error {
	h.sugar.Info("PostArticlesDelete")
	articleID := ectx.Param("id")
	_, userID, err := h.parseSession(ectx)
	if err != nil {
		return err
	}
	article, err := h.articleStore.Get(ectx.Request().Context(), articleID)
	if err != nil {
		return err
	}
	if userID != article.UserID {
		return ectx.String(http.StatusForbidden, "you are not the owner of this article")
	}
	if err := h.articleStore.Delete(ectx.Request().Context(), articleID); err != nil {
		return err
	}
	return ectx.Redirect(http.StatusFound, "/auth/articles")
}

func (h *handlerImpl) LoginRequired(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ectx echo.Context) error {
		h.sugar.Info("LoginRequired middleware")
		_, userID, err := h.parseSession(ectx)
		if err != nil {
			h.sugar.Infow("session does not exist", "cause", err)
			return ectx.Redirect(http.StatusFound, "/login")
		}
		h.sugar.Infow("session exists", "userID", userID)
		return next(ectx)
	}
}

func (h *handlerImpl) parseSession(ectx echo.Context) (string, string, error) {
	sess, err := session.Get(constSessionName, ectx)
	if err != nil {
		return "", "", fmt.Errorf("session cookie does not exist")
	}
	log.Printf("sess.Values: %v\n", sess.Values)
	sessionID := sess.Values[constSessionID].(string)
	sessionStruct, err := h.sessionStore.Get(ectx.Request().Context(), sessionID)
	if err != nil {
		return "", "", fmt.Errorf("session %s does not exist", sessionID)
	}
	return sessionID, sessionStruct.UserID, nil
}
