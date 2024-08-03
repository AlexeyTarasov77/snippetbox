package forms

type SnippetCreateForm struct {
	Title string `schema:"title, required" validate:"required,max=100"`
	Content string `schema:"content, required" validate:"required"`
	Expires int `schema:"expires, required" validate:"required,gte=1,lte=365"`
	BaseForm
}

