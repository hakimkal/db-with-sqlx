package main

import (
	"context"
	"fmt"
	"log"

	"github.com/hakimkal/db-with-sqlx/internal/config"
	"github.com/hakimkal/db-with-sqlx/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"

	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func main() {

	var ctx = context.Background()
	var err error

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	pool, err := pgxpool.New(ctx, cfg.DbUrl)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	dbPool := stdlib.OpenDBFromPool(pool)
	defer pool.Close()

	if err = pool.Ping(ctx); err != nil {
		log.Fatalf("Database ping failed: %v\n", err)
	}
	db = sqlx.NewDb(dbPool, "pgx")
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	defer db.Close()

	fmt.Println("Database successfully connected and initialized.")

	var userService = service.DbService{Db: db}

	users, err := userService.ListUsers()
	if err != nil {
		log.Printf("%v", err)
	}
	for _, user := range users {

		fmt.Printf("%d | %s | %s \n", user.Id, user.Name, user.Email)
	}
	//log.Printf("%+v\n", users)
}
