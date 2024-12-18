package handler

import (
	"html/template"
	"io"
)

type SpecialTemplate struct {
	Template *template.Template
	Name     string
}

func ParsePage(name string) *SpecialTemplate {
	return &SpecialTemplate{Template: template.Must(template.ParseFiles("templates/templates.gohtml", "templates/"+name+".gohtml")),
		Name: name + ".gohtml"}
}

func (t *SpecialTemplate) Execute(w io.Writer, data any) error {
	return t.Template.ExecuteTemplate(w, t.Name, data)
}
