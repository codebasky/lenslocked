package controllers

import (
	"fmt"
	"net/http"
	"net/url"

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
	esrv       *model.EmailService
	psrv       *model.PasswordResetService
}

func New(in Template, up Template,
	usrv *model.UserService, ssrv *model.SessionService,
	esrv *model.EmailService, psrv *model.PasswordResetService) *User {
	return &User{
		signinTmpl: in,
		signupTmpl: up,
		usrv:       usrv,
		ssrv:       ssrv,
		esrv:       esrv,
		psrv:       psrv,
	}
}

func (u User) Signin(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	fmt.Println(data.Email)
	u.signinTmpl.Execute(w, r, data)
}

func (u User) ProcessSignIn(w http.ResponseWriter, r *http.Request) {
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
		return
	}
	setCookie(w, SessionCookie, session.Token)
	http.Redirect(w, r, "/users/me", http.StatusFound)
}

func (u User) Signup(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.signupTmpl.Execute(w, r, data)
}

func (u User) ProcessSignup(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	//fmt.Fprintf(w, "user type email %s pwd %s", email, password)
	user, err := u.usrv.Create(email, password)
	if err != nil {
		fmt.Printf("User creation failed with error: %s", err)
		http.Error(w, "User Signup failed", http.StatusInternalServerError)
		return
	}

	session, err := u.ssrv.Create(user.ID)
	if err != nil {
		fmt.Printf("User session creation failed with error: %s", err)
		http.Error(w, "User Signup failed", http.StatusInternalServerError)
		return
	}
	setCookie(w, SessionCookie, session.Token)
	http.Redirect(w, r, "/users/me", http.StatusFound)
}

func (u User) ProcessSignout(w http.ResponseWriter, r *http.Request) {
	token, err := readCookie(r, SessionCookie)
	if err != nil {
		fmt.Printf("User session token missing: %s", err)
		http.Error(w, "User Signout failed", http.StatusInternalServerError)
		return
	}
	//fmt.Fprintf(w, "user type email %s pwd %s", email, password)
	err = u.ssrv.Delete(token)
	if err != nil {
		fmt.Printf("User session deletion failed: %s", err)
		http.Error(w, "User Signout failed", http.StatusInternalServerError)
		return
	}

	deleteCookie(w, SessionCookie)
	http.Redirect(w, r, "/signin", http.StatusFound)
}

func (u User) CurrentUser(w http.ResponseWriter, r *http.Request) {
	token, err := readCookie(r, SessionCookie)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
	user, err := u.ssrv.User(token)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
	fmt.Fprintf(w, "Current user: %s\n", user.Email)
}

func (u User) ProcessForgotPwd(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	user, err := u.usrv.User(email)
	if err != nil {
		fmt.Printf("User find failed with error: %s", err)
		http.Error(w, "reset password failed wrong emailid", http.StatusBadRequest)
		return
	}

	pwReset, err := u.psrv.Create(user.Email, user.ID)
	if err != nil {
		fmt.Printf("Process forgot pwd failed: %s", err)
		http.Error(w, "pwd create failed", http.StatusInternalServerError)
		return
	}
	vals := url.Values{
		"token": {pwReset.Token},
	}
	resetURL := "https://www.lenslocked.com/reset-pw?" + vals.Encode()

	err = u.esrv.ForgotPassword(user.Email, resetURL)
	if err != nil {
		fmt.Printf("Process forgot pwd failed: %s", err)
		http.Error(w, "reset password failed", http.StatusInternalServerError)
		return
	}
}
