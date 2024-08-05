package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"snippetbox.proj.net/internal/api/constants"
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
	snippetID, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.logger.Error("Error inserting snippet", "err", err.Error())
		response.HttpError(w, "", http.StatusInternalServerError)
		return
	}
	app.sessionManager.Put(r.Context(), "flash", "Snippet successfully created!")
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", snippetID), http.StatusSeeOther)
}

func (app *Application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = forms.SnippetCreateForm{
		Expires: 1,
	}
	app.render(w, "create.html", data)

}

func (app *Application) snippetView(w http.ResponseWriter, r *http.Request) {
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

func (app *Application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = forms.UserSignupForm{}
	app.render(w, "signup.html", data)
}

func (app *Application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	rerenderTemplate := func(form forms.UserSignupForm) {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, "signup.html", data, http.StatusUnprocessableEntity)
	}
	if err := r.ParseForm(); err != nil {
		app.logger.Error("Error parsing form", "err", err.Error())
		response.HttpError(w, "", http.StatusBadRequest)
		return
	}
	var form forms.UserSignupForm
	if err := app.formDecoder.Decode(&form, r.PostForm); err != nil {
		app.logger.Error("Error decoding form", "err", err.Error())
		response.HttpError(w, "", http.StatusBadRequest)
		return
	}
	if !form.IsValid(form) {
		rerenderTemplate(form)
		return
	}
	id, err := app.users.Insert(form.Username, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, storage.ErrDuplicateEmail) {
			form.FieldErrors["email"] = "Address is already in use"
			rerenderTemplate(form)
			return
		}
		app.logger.Error("Error inserting user", "err", err.Error())
		response.HttpError(w, "")
		return
	}
	app.logger.Info("Created new user. Redirecting to login page", "id", id, "formData", form)
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *Application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	app.sessionManager.Remove(r.Context(), string(constants.UserIDCtxKey))
	app.sessionManager.Put(r.Context(), string(constants.FlashCtxKey), "You've been logged out successfully!")
	app.sessionManager.RenewToken(r.Context())
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *Application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = forms.UserLoginForm{}
	app.render(w, "login.html", data)
}

func (app *Application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		app.logger.Error("Error parsing form", "err", err.Error())
		response.HttpError(w, "", http.StatusBadRequest)
		return
	}
	app.logger.Debug("Form parsed", "data", r.PostForm)
	var form forms.UserLoginForm
	if err := app.formDecoder.Decode(&form, r.PostForm); err != nil {
		app.logger.Error("Error decoding form", "err", err.Error())
		response.HttpError(w, "", http.StatusBadRequest)
		return
	}
	if !form.IsValid(form) {
		app.logger.Debug("Form is not valid", "formErrors", form.FieldErrors, "formData", form)
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, "login.html", data, http.StatusUnprocessableEntity)
		return
	}
	user, err := app.users.Authenticate(form.Email, form.Password)
	if errors.Is(err, storage.ErrInvalidCredentials) {
		app.logger.Info("Invalid credentials provided", "err", err.Error(), "formData", form)
		form.NonFieldErrors = append(form.NonFieldErrors, "Invalid email or password")
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, "login.html", data, http.StatusUnprocessableEntity)
		return
	}
	app.sessionManager.RenewToken(r.Context())
	app.sessionManager.Put(
		r.Context(),
		string(constants.FlashCtxKey),
		fmt.Sprintf("Hello, %s! You've login succesfully", user.Username),
	)
	app.sessionManager.Put(r.Context(), string(constants.UserIDCtxKey), user.ID)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
