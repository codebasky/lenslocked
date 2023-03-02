package controllers

import (
	"fmt"
	"net/http"

	"github.com/codebasky/lenslocked/model"
)

type User struct {
	signinTmpl Template
	signupTmpl Template
	srv        *model.UserService
}

func New(in Template, up Template, usrv *model.UserService) *User {
	return &User{
		signinTmpl: in,
		signupTmpl: up,
		srv:        usrv,
	}
}

func (u User) Signin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data struct {
			Email string
		}
		data.Email = r.FormValue("email")
		fmt.Println(data.Email)
		u.signinTmpl.Execute(w, r, data)
	}
}

func (u User) ProcessSignIn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := r.FormValue("email")
		password := r.FormValue("password")
		//fmt.Fprintf(w, "user type email %s pwd %s", email, password)
		status, err := u.srv.Authenticate(email, password)
		if err != nil || !status {
			fmt.Printf("User sigup failed with error: %s", err)
			http.Error(w, "User Signup failed", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "User Authenticated")
	}
}

func (u User) Signup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data struct {
			Email string
		}
		data.Email = r.FormValue("email")
		u.signupTmpl.Execute(w, r, data)
	}
}

func (u User) ProcessSignup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := r.FormValue("email")
		password := r.FormValue("password")
		//fmt.Fprintf(w, "user type email %s pwd %s", email, password)
		user, err := u.srv.Create(email, password)
		if err != nil {
			fmt.Printf("User creation failed with error: %s", err)
			http.Error(w, "User Signup failed", http.StatusInternalServerError)
		}
		fmt.Fprintf(w, "User created: %+v", *user)
	}
}
