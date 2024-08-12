package templates

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"time"

	"snippetbox.proj.net/internal/api/forms"
	"snippetbox.proj.net/internal/storage/models"
	"snippetbox.proj.net/ui"
)

type TemplateData struct {
	CurrentYear int
	Snippet 	*models.Snippet
	Snippets 	[]*models.Snippet
	Form        any
	Flash       string
	User        *models.User
	CSRFToken 	string
}

// func (tmplData *TemplateData) IsRequiredField(fieldName string) bool {
// 	field, found := reflect.TypeOf(tmplData.Form).FieldByName(fieldName)
// 	if !found {
// 		panic(fmt.Sprintf("Field %s not found in type %s", fieldName, reflect.TypeOf(tmplData.Form).Name()))
// 	}
// 	if rules := field.Tag.Get("validate"); rules != "" {
// 		if strings.Contains(rules, "required") {
// 			return true
// 		}
// 	}
// 	return false
// }

func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format("02 Jan 2006 at 15:04")
}

var funcMap = template.FuncMap{
	"humanDate": humanDate,
	"isRequiredFormField": forms.IsRequiredField,
}

func NewTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	pages, err := fs.Glob(ui.Files, "html/pages/*.html")
	if err != nil {
		return nil, err
	}
	partials, err := fs.Glob(ui.Files, "html/partials/*.html")
	if err != nil {
		return nil, err
	}

	for _, path := range pages {
		filename := filepath.Base(path)
		relevantFiles := append([]string{"html/base.html"}, partials...)
		relevantFiles = append(relevantFiles, path)
		parsed, err := template.New(filename).Funcs(funcMap).ParseFS(ui.Files, relevantFiles...)
		if err != nil {
			return nil, err
		}
		cache[filename] = parsed
	}
	return cache, nil
}
