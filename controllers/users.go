package controllers

import (
	"fmt"
	"net/http"

	"github.com/codebasky/lenslocked/model"
)

const (
	SessionCookie = "session"
)

type User struct {
	signinTmpl Template
	signupTmpl Template
	usrv       *model.UserService
	ssrv       *model.SessionService
}

func New(in Template, up Template, usrv *model.UserService, ssrv *model.SessionService) *User {
	return &User{
		signinTmpl: in,
		signupTmpl: up,
		usrv:       usrv,
		ssrv:       ssrv,
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
		user, err := u.usrv.Authenticate(email, password)
		if err != nil {
			fmt.Printf("User sigup failed with error: %s", err)
			http.Error(w, "User Signup failed", http.StatusInternalServerError)
			return
		}
		session, err := u.ssrv.Create(user.ID)
		if err != nil {
			fmt.Printf("User session creation failed with error: %s", err)
			http.Error(w, "User Signin failed", http.StatusInternalServerError)
		}
		setCookie(w, SessionCookie, session.Token)
		fmt.Fprintf(w, "User Signin success")
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
		user, err := u.usrv.Create(email, password)
		if err != nil {
			fmt.Printf("User creation failed with error: %s", err)
			http.Error(w, "User Signup failed", http.StatusInternalServerError)
		}

		session, err := u.ssrv.Create(user.ID)
		if err != nil {
			fmt.Printf("User session creation failed with error: %s", err)
			http.Error(w, "User Signup failed", http.StatusInternalServerError)
		}
		setCookie(w, SessionCookie, session.Token)
		fmt.Fprintf(w, "User created: %+v", *user)
	}
}
