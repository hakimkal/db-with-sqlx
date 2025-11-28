package model

type User struct {
	Id    int    `db:"id"`
	Name  string `db:"name"`
	Email string `db:"email"`
}
