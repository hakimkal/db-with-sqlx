package service

import (
	"github.com/hakimkal/db-with-sqlx/internal/model"
	"github.com/jmoiron/sqlx"
)

type User model.User
type UserService interface {
	GetUser(Id int) (*User, error)
	ListUsers() ([]User, error)
	CreateUser(newUser model.User) (*User, error)
}

type DbService struct {
	Db *sqlx.DB
}

func (s *DbService) CreateUser(newUser User) (*User, error) {

	query := "INSERT INTO users (name, email) values ( :name, :email)" +
		"RETURNING * "
	err := sqlx.Get(s.Db, &newUser, query, newUser)
	return &newUser, err
}

func (s *DbService) GetUser(Id int) (*User, error) {
	var user User
	err := s.Db.Get(&user, "SELECT id, name, email FROM users WHERE id = $1", Id)
	return &user, err
}

func (s *DbService) ListUsers() ([]User, error) {
	var users []User
	err := s.Db.Select(&users, "SELECT id, name, email FROM users ORDER BY  id ASC ")
	return users, err
}
