package views

import (
	"log"
	"net/http"
	"path/filepath"
	"text/template"
)

type Template struct {
	httpTpl *template.Template
}

func Must(t Template, err error) Template {
	if err != nil {
		panic(err)
	}
	return t
}

func Parse(fname string) (Template, error) {
	tpl, err := template.ParseFiles(filepath.Join("templates", fname))
	if err != nil {
		return Template{}, err
	}
	return Template{
		httpTpl: tpl,
	}, nil
}

func (t Template) Execute(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := t.httpTpl.Execute(w, data)
	if err != nil {
		log.Printf("executing template: %v", err)
		http.Error(w, "There was an error executing the template.", http.StatusInternalServerError)
		return
	}
}
