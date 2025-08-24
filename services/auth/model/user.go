package model

type User struct {
	ID       int    `json:"id"`
	EMail    string `json:"email"`
	Password string `json:"passwod"`
}
