package models

type User struct {
	ID           int    `db:"id" json:"id"`
	Login        string `db:"login" json:"login"`
	PasswordHash string `db:"password_hash" json:"-"`
}
