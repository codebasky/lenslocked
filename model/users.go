package model

import (
	"database/sql"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID            int
	Email         string
	Password_Hash string
}

type UserService struct {
	db *sql.DB
}

func NewUserSrv(db *sql.DB) *UserService {
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
	err = result.Scan(&user.ID)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u UserService) Authenticate(email string, pwd string) (*User, error) {
	emailId := strings.ToLower(email)

	result := u.db.QueryRow(`select id, password_hash from users where email=$1`, emailId)
	user := User{
		Email: emailId,
	}
	err := result.Scan(&user.ID, &user.Password_Hash)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password_Hash), []byte(pwd))
	if err != nil {
		return nil, err
	}
	return &user, nil
}
