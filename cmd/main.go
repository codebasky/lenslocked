package main

import (
	"fmt"
	"net/http"

	"github.com/codebasky/lenslocked/controllers"
	"github.com/codebasky/lenslocked/templates"
	"github.com/codebasky/lenslocked/views"
)

func main() {
	fmt.Println("Starting web development")
	http.HandleFunc("/", controllers.StaticHandler(views.Must(views.ParseFS(templates.FS, "home.gohtml", "tailwind.gohtml"))))
	http.HandleFunc("/contact", controllers.StaticHandler(views.Must(views.ParseFS(templates.FS, "contact.gohtml", "tailwind.gohtml"))))
	http.HandleFunc("/faq", controllers.FAQ(views.Must(views.ParseFS(templates.FS, "faq.gohtml", "tailwind.gohtml"))))
	http.ListenAndServe(":3000", nil)
}
