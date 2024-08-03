package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"snippetbox.proj.net/internal/api/forms"
	"snippetbox.proj.net/internal/api/response"
	"snippetbox.proj.net/internal/storage"
)

func (app *Application) home(w http.ResponseWriter, r *http.Request) {
	latestSnippets, err := app.snippets.Latest(10)
	if err != nil && !errors.Is(err, storage.ErrNoRecord) {
		response := response.Error("Database error")
		http.Error(w, response.Error, response.Status)
		return
	}
	data := app.newTemplateData(r)
	data.Snippets = latestSnippets
	app.logger.Debug("home", "snippets", latestSnippets, "data", data)
	app.render(w, "home.html", data)
}

func (app *Application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		app.logger.Error("Error parsing form", "err", err.Error())
		response.HttpError(w, "", http.StatusBadRequest)
		return
	}
	var form forms.SnippetCreateForm
	if err := app.formDecoder.Decode(&form, r.PostForm); err != nil {
		app.logger.Error("Error decoding form", "err", err.Error())
		response.HttpError(w, "", http.StatusBadRequest)
		return
	}
	if !form.IsValid(form) {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, "create.html", data, http.StatusUnprocessableEntity)
		return
	}
	snippetID, err := app.snippets.Insert(form.Title, form.Content, form.Expires);
	if err != nil {
		app.logger.Error("Error inserting snippet", "err", err.Error())
		response.HttpError(w, "", http.StatusInternalServerError)
		return
	}
	app.sessionManager.Put(r.Context(), "flash", "Snippet successfully created!")
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", snippetID), http.StatusSeeOther)
}

func (app *Application) snippetCreateGet(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = forms.SnippetCreateForm{
		Expires: 1,
	}
	app.render(w, "create.html", data)
	
}

func (app *Application) snippetView(w http.ResponseWriter, r *http.Request,) {
	id := chi.URLParam(r, "id")
	idInt, err := strconv.Atoi(id)
	app.logger.Debug("Got value", "id", id)
	if err != nil || idInt < 1 {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	snippet, err := app.snippets.Get(idInt)
	if err != nil {
		if errors.Is(err, storage.ErrNoRecord) {
			http.Error(w, "Snippet not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}
	data := app.newTemplateData(r)
	data.Snippet = snippet
	app.render(w, "view.html", data)
}

func (app *Application) snippetNotFound(w http.ResponseWriter, r *http.Request) {
	app.render(w, "404.html", app.newTemplateData(r))
}
