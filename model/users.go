package model

import (
	"database/sql"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id            int
	Email         string
	Password_Hash string
}

type UserService struct {
	db *sql.DB
}

func New(db *sql.DB) *UserService {
	return &UserService{
		db: db,
	}
}

func (u UserService) Create(email string, pwd string) (*User, error) {
	emailId := strings.ToLower(email)
	pHash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	result := u.db.QueryRow(`insert into users (email, password_hash) 
							values ($1, $2) returning id`, emailId, pHash)
	user := User{
		Email:         emailId,
		Password_Hash: string(pHash),
	}
	err = result.Scan(&user.Id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u UserService) Authenticate(email string, pwd string) (bool, error) {
	emailId := strings.ToLower(email)

	result := u.db.QueryRow(`select password_hash from users where email=$1`, emailId)
	var pHash string
	err := result.Scan(&pHash)
	if err != nil {
		return false, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(pHash), []byte(pwd))
	if err != nil {
		return false, err
	}
	return true, nil
}
