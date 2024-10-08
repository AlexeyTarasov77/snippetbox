package main

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/justinas/nosurf"
	"snippetbox.proj.net/internal/api/constants"
	"snippetbox.proj.net/internal/api/response"
	"snippetbox.proj.net/internal/storage/models"
	"snippetbox.proj.net/internal/templates"
)

func (app *Application) render(
	w http.ResponseWriter, name string,
	data *templates.TemplateData, _status ...int,
) {
	if len(_status) == 0 {
		_status = append(_status, http.StatusOK)
	}
	status := _status[0]
	ts, exists := app.templateCache[name]
	if !exists {
		panic(fmt.Sprintf("The template %s does not exist", name))
	}
	buffer := new(bytes.Buffer)
	if err := ts.ExecuteTemplate(buffer, "base", data); err != nil {
		app.logger.Error("Error executing template", "name", name, "err", err.Error())
		response.HttpError(w, "Error to render template")
		return
	}
	w.WriteHeader(status)
	buffer.WriteTo(w)
}

func (app *Application) rerenderInvalidForm(
	w http.ResponseWriter,
	r *http.Request,
	form any,
	templateName string,
) {
	data := app.newTemplateData(r)
	data.Form = form
	app.render(w, templateName, data, http.StatusUnprocessableEntity)
}

func (app *Application) newTemplateData(r *http.Request) *templates.TemplateData {
	user, ok := r.Context().Value(constants.UserCtxKey).(*models.User)
	if !ok {
		user = nil
	}
	return &templates.TemplateData{
		CurrentYear: time.Now().Year(),
		Flash: app.sessionManager.PopString(r.Context(), "flash"),
		User:  user,
		CSRFToken: nosurf.Token(r),
	}
}
