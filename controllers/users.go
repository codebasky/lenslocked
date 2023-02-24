package controllers

import (
	"fmt"
	"net/http"
)

type User struct {
	tmpl Template
}

func New(t Template) *User {
	return &User{
		tmpl: t,
	}
}

func (u User) Signup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data struct {
			Email string
		}
		data.Email = r.FormValue("email")
		fmt.Println(data.Email)
		u.tmpl.Execute(w, data)
	}
}

func (u User) ProcessSignup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := r.FormValue("email")
		password := r.FormValue("password")
		fmt.Fprintf(w, "user type email %s pwd %s", email, password)
	}
}
