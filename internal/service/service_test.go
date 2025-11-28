package service

import (
	"context"

	"log"
	"testing"

	"github.com/hakimkal/db-with-sqlx/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"

	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

// Helper function
func setupTestDB(t *testing.T) *DbService {
	var ctx = context.Background()
	var err error
	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatal(err)
	}

	pool, err := pgxpool.New(ctx, cfg.TestDbUrl)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	dbPool := stdlib.OpenDBFromPool(pool)
	//defer pool.Close()

	if err = pool.Ping(ctx); err != nil {
		log.Fatalf("Database ping failed: %v\n", err)
	}
	db = sqlx.NewDb(dbPool, "pgx")
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	//defer db.Close()

	// Set up necessary tables for the test if they don't exist

	createTableQuery := `CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        name VARCHAR(100),
        email VARCHAR(100),
        active BOOLEAN DEFAULT TRUE
    );`
	if _, err := db.ExecContext(ctx, createTableQuery); err != nil {
		log.Fatalf("Failed to prepare test table: %v", err)
	}

	return &DbService{Db: db}
}

func TestDBService_ListUsers(t *testing.T) {
	// Arrange: Use the real test DB setup
	service := setupTestDB(t)
	ctx := context.Background()

	// Arrange: Insert some test data directly into the DB for the test to read
	insertQuery := "INSERT INTO users (name, email) VALUES ($1, $2)"
	service.Db.ExecContext(ctx, insertQuery, "Test User 1", "test1@example.com")
	service.Db.ExecContext(ctx, insertQuery, "Test User 2", "test2@example.com")

	// Act: Call the real method on the real DB service
	users, err := service.ListUsers()

	// Assert: Verify the results
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(users) != 2 {
		t.Errorf("expected 2 users, got %d", len(users))
	}

	if users[0].Name != "Test User 1" {
		t.Errorf("expected user name 'Test User 1', got %s", users[0].Name)
	}

	// Clean up specific data after the test
	service.Db.ExecContext(ctx, "DELETE FROM users")
	// Ensure we clean up after the test suite finishes
	t.Cleanup(func() {
		db.Close()
	})

}
