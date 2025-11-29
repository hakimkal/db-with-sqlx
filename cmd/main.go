package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

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

	taskName := flag.String("taskName", "none", "Specify the task: 'list' or 'view-user'")

	userID := flag.Int("id", 0, "The ID of the user to view (required for view-users task)")

	flag.Parse()

	switch *taskName {

	case "list":
		var userService = service.DbService{Db: db}

		users, err := userService.ListUsers()
		if err != nil {
			log.Printf("%v", err)
		}
		for _, user := range users {

			fmt.Printf("%d | %s | %s \n", user.Id, user.Name, user.Email)
		}
	case "view":
		if *userID <= 0 {
			fmt.Println("Error: The 'view-users' task requires a valid positive -id flag.")
			flag.Usage() // Show user how to use the flags
			os.Exit(1)
		}
		var userService = service.DbService{Db: db}
		user, err := userService.GetUser(*userID)
		if err != nil {
			log.Printf("Select user error | %v", err)

		}
		if user == nil {
			fmt.Println("User not found")
		} else {
			fmt.Printf("%s | %s | %s \n", user.Id, user.Name, user.Email)

		}

	case "none":
		// This runs if the user didn't provide the flag, or used the default.
		fmt.Println("You need to specify a taskName.")
		flag.Usage() // Prints usage instructions

	default:

		fmt.Printf("Unknown task specified: %s\n", *taskName)
		flag.Usage()
	}

}
