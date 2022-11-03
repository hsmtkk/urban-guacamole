package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/hsmtkk/urban-guacamole/entry"
	"github.com/hsmtkk/urban-guacamole/entryrepo"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type MyController struct {
	sugar     *zap.SugaredLogger
	entryRepo entryrepo.EntryRepo
}

func New(sugar *zap.SugaredLogger, entryRepo entryrepo.EntryRepo) *MyController {
	return &MyController{sugar, entryRepo}
}

func (ctrl *MyController) ShowEntries(ectx echo.Context) error {
	ctrl.sugar.Info("ShowEntries")
	entries, err := ctrl.entryRepo.Scan()
	if err != nil {
		return err
	}
	type entry struct {
		Title        string
		ShowEntryURL string
	}
	type tmplParam struct {
		Entries []entry
	}
	var param tmplParam
	for _, e := range entries {
		param.Entries = append(param.Entries, entry{
			Title:        e.Title,
			ShowEntryURL: fmt.Sprintf("/entries/%d", e.ID),
		})
	}
	return ectx.Render(http.StatusOK, "index", param)
}

func (ctrl *MyController) AddEntry(ectx echo.Context) error {
	ctrl.sugar.Info("AddEntry")
	title := ectx.FormValue("title")
	text := ectx.FormValue("text")
	e := entry.New(title, text)
	if err := ctrl.entryRepo.Save(e); err != nil {
		return err
	}
	return ectx.Redirect(http.StatusSeeOther, "/")
}

func (ctrl *MyController) NewEntry(ectx echo.Context) error {
	ctrl.sugar.Info("NewEntry")
	return ectx.Render(http.StatusOK, "new", nil)
}

func (ctrl *MyController) ShowEntry(ectx echo.Context) error {
	ctrl.sugar.Info("ShowEntry")
	idStr := ectx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse %s as int64; %w", idStr, err)
	}
	e, err := ctrl.entryRepo.Get(id)
	if err != nil {
		return err
	}
	type tmplParam struct {
		Title          string
		Text           string
		CreatedAt      string
		EditEntryURL   string
		DeleteEntryURL string
	}
	param := tmplParam{
		Title:          e.Title,
		Text:           e.Text,
		CreatedAt:      e.CreatedAt.Format("2006-01-02 15:04:05"),
		EditEntryURL:   fmt.Sprintf("/entries/%d/edit", e.ID),
		DeleteEntryURL: fmt.Sprintf("/entries/%d/delete", e.ID),
	}
	return ectx.Render(http.StatusOK, "show", param)
}

func (ctrl *MyController) EditEntry(ectx echo.Context) error {
	ctrl.sugar.Info("EditEntry")
	idStr := ectx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse %s as int64; %w", idStr, err)
	}
	e, err := ctrl.entryRepo.Get(id)
	if err != nil {
		return err
	}
	type tmplParam struct {
		Title          string
		Text           string
		UpdateEntryURL string
	}
	param := tmplParam{
		Title:          e.Title,
		Text:           e.Text,
		UpdateEntryURL: fmt.Sprintf("/entries/%d/update", e.ID),
	}
	return ectx.Render(http.StatusOK, "edit", param)
}

func (ctrl *MyController) UpdateEntry(ectx echo.Context) error {
	ctrl.sugar.Info("UpdateEntry")
	idStr := ectx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse %s as int64; %w", idStr, err)
	}
	e, err := ctrl.entryRepo.Get(id)
	if err != nil {
		return err
	}
	e.Title = ectx.FormValue("title")
	e.Text = ectx.FormValue("text")
	if err := ctrl.entryRepo.Save(e); err != nil {
		return err
	}
	return ectx.Redirect(http.StatusSeeOther, "/")
}

func (ctrl *MyController) DeleteEntry(ectx echo.Context) error {
	ctrl.sugar.Info("DeleteEntry")
	idStr := ectx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse %s as int64; %w", idStr, err)
	}
	if err := ctrl.entryRepo.Delete(id); err != nil {
		return err
	}
	return ectx.Redirect(http.StatusSeeOther, "/")
}
