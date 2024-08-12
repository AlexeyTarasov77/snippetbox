package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"
	"snippetbox.proj.net/internal/api/constants"
	"snippetbox.proj.net/internal/api/forms"
	"snippetbox.proj.net/internal/api/response"
	"snippetbox.proj.net/internal/storage"
	"snippetbox.proj.net/internal/storage/models"
)


func (app *Application) home(w http.ResponseWriter, r *http.Request) {
	latestSnippets, err := app.snippets.Latest(10)
	if err != nil && !errors.Is(err, storage.ErrNoRecord) {
		app.logger.Error("Error getting latest snippets", "err", err.Error())
		response.HttpError(w, "Database error")
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
	currUserID := app.sessionManager.GetInt(r.Context(), string(constants.UserIDCtxKey))
	snippetID, err := app.snippets.Insert(form.Title, form.Content, form.Expires, currUserID)
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
			response.HttpError(w, "", http.StatusNotFound)
		} else {
			response.HttpError(w, "")
		}
		return
	}
	data := app.newTemplateData(r)
	data.Snippet = snippet
	app.render(w, "view.html", data)
}

func (app *Application) snippetNotFound(w http.ResponseWriter, r *http.Request) {
	app.logger.Info("Not found", "path", r.URL.Path)
	app.render(w, "404.html", app.newTemplateData(r))
}

func (app *Application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = forms.UserSignupForm{}
	app.render(w, "signup.html", data)
}

func (app *Application) userSignupPost(w http.ResponseWriter, r *http.Request) {
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
		app.rerenderInvalidForm(w, r, form, "signup.html")
		return
	}
	id, err := app.users.Insert(form.Username, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, storage.ErrDuplicateEmail) {
			app.logger.Info("Duplicate email", "email", form.Email)
			form.FieldErrors["email"] = "Address is already in use"
			app.rerenderInvalidForm(w, r, form, "signup.html")
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
		app.rerenderInvalidForm(w, r, form, "login.html")
		return
	}
	user, err := app.users.Authenticate(form.Email, form.Password)
	if errors.Is(err, storage.ErrInvalidCredentials) {
		app.logger.Info("Invalid credentials provided", "err", err.Error(), "formData", form)
		form.NonFieldErrors = append(form.NonFieldErrors, "Invalid email or password")
		app.rerenderInvalidForm(w, r, form, "login.html")
		return
	}
	app.sessionManager.RenewToken(r.Context())
	app.sessionManager.Put(
		r.Context(),
		string(constants.FlashCtxKey),
		fmt.Sprintf("Hello, %s! You've login succesfully", user.Username),
	)
	app.sessionManager.Put(r.Context(), string(constants.UserIDCtxKey), user.ID)
	redirectPath := app.sessionManager.PopString(r.Context(), string(constants.RedirectCtxKey))
	if redirectPath == "" {
		redirectPath = "/user/account"
	}
	http.Redirect(w, r, redirectPath, http.StatusSeeOther)
}

func (app *Application) accountView(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	userSnippets, err := app.snippets.GetByUserID(data.User.ID)
	if err != nil && !errors.Is(err, storage.ErrNoRecord) {
		app.logger.Error("Error getting user snippets", "err", err.Error())
		response.HttpError(w, "")
		return
	}
	data.Snippets = userSnippets
	app.render(w, "account.html", data)
}

func (app *Application) accountPasswordUpdate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = forms.UserPasswordUpdateForm{}
	app.render(w, "password-update.html", data)
}

func (app *Application) accountPasswordUpdatePost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		app.logger.Error("Error parsing form", "err", err.Error())
		response.HttpError(w, "", http.StatusBadRequest)
		return
	}
	app.logger.Debug("Form parsed", "data", r.PostForm)
	var form forms.UserPasswordUpdateForm
	if err := app.formDecoder.Decode(&form, r.PostForm); err != nil {
		app.logger.Error("Error decoding form", "err", err.Error())
		response.HttpError(w, "", http.StatusBadRequest)
		return
	}
	if !form.IsValid(form) {
		app.logger.Debug("Form is not valid", "formData", form)
		app.rerenderInvalidForm(w, r, form, "password-update.html")
		return
	}
	currUser := r.Context().Value(constants.UserCtxKey).(*models.User)
	_, err := app.users.Authenticate(currUser.Email, form.CurrentPassword)
	if errors.Is(err, storage.ErrInvalidCredentials) {
		app.logger.Info("Invalid credentials provided", "err", err.Error(), "formData", form)
		form.NonFieldErrors = append(form.NonFieldErrors, "Invalid credentials")
		app.rerenderInvalidForm(w, r, form, "password-update.html")
		return
	}
	newPasswordHash, err := bcrypt.GenerateFromPassword([]byte(form.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		app.logger.Error("Error hashing password", "err", err.Error())
		response.HttpError(w, "")
		return
	}
	currUser.Password = newPasswordHash
	if err := app.users.Update(currUser); err != nil {
		app.logger.Error("Error updating user", "err", err.Error())
		response.HttpError(w, "")
		return
	}
	app.sessionManager.Put(r.Context(), string(constants.FlashCtxKey), "Password successfully updated!")
	http.Redirect(w, r, "/user/account", http.StatusSeeOther)
}

func (app *Application) about(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, "about.html", data)
}

func (app *Application) ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
