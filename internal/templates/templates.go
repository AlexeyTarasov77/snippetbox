package templates

import (
	"html/template"
	"path/filepath"
	"time"

	"snippetbox.proj.net/internal/storage/models"
)

type TemplateData struct {
	CurrentYear int
	Snippet 	*models.Snippet
	Snippets 	[]*models.Snippet
	Form        any
	Flash       string
	User        *models.User
}

// func (tmplData *TemplateData) IsRequiredField(fieldName string) bool {
// 	field, found := reflect.TypeOf(tmplData.Form).FieldByName(fieldName)
// 	if !found {
// 		panic(fmt.Sprintf("Field %s not found in type %s", fieldName, reflect.TypeOf(tmplData.Form).Name()))
// 	}
// 	return field.Tag.Get("required") != ""
// }

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var funcMap = template.FuncMap{
	"humanDate": humanDate,
}

func NewTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	pages, err := filepath.Glob("ui/html/pages/*.html")
	if err != nil {
		return nil, err
	}
	partials, err := filepath.Glob("ui/html/partials/*.html")
	if err != nil {
		return nil, err
	}

	for _, path := range pages {
		filename := filepath.Base(path)
		relevantFiles := append([]string{"ui/html/base.html"}, partials...)
		relevantFiles = append(relevantFiles, path)
		parsed, err := template.New(filename).Funcs(funcMap).ParseFiles(relevantFiles...)
		if err != nil {
			return nil, err
		}
		cache[filename] = parsed
	}
	return cache, nil
}


// htmlfiles paths should be specified relative to ui/html dir
// func RenderTemplate(w http.ResponseWriter, data any, htmlfiles ...string) {
// 	baseFileName := strings.Replace(htmlfiles[0], ".html", "", 1)
// 	for i := range htmlfiles {
// 		htmlfiles[i] = "ui/html/" + htmlfiles[i]
// 	}
// 	parsed, err := template.ParseFiles(htmlfiles...)
// 	if err != nil {
// 		slog.Error("Error parsing template", "msg", err)
// 		response := response.Error("Error parsing template")
// 		http.Error(w, response.Error, response.Status)
// 		return
// 	}
// 	parsed.ExecuteTemplate(w, baseFileName, data)
// }


// Shortcut for render template including base and partials
// func RenderPage(w http.ResponseWriter, data any, htmlfile string) {
// 	RenderTemplate(w, data, "base.html", htmlfile, "partials/nav.html")
// }
