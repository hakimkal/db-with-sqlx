package config

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
)

type Config struct {
	DbUrl     string
	TestDbUrl string
}

// init is a good place to load environment variables for the entire package lifecycle
func init() {
	// Find the project root by searching for go.mod file
	projectRoot := findProjectRoot()
	if projectRoot == "" {
		log.Fatal("Could not find project root (go.mod file) to load .env")
	}

	envPath := filepath.Join(projectRoot, ".env")

	// Load the .env file from the determined path
	if err := godotenv.Load(envPath); err != nil {
		 
		log.Fatalf("Error loading .env file from %s: %v", envPath, err)
	}
}

func LoadConfig() (*Config, error) {

	config := &Config{
		DbUrl:     os.Getenv("DB_URL"),
		TestDbUrl: os.Getenv("TEST_DB_URL"),
	}
	return config, nil
}

// Helper function to find the project root by locating the go.mod file
func findProjectRoot() string {
	// Start searching from the directory where this source file lives
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return ""
	}
	currentDir := filepath.Dir(filename)

	for {
		if _, err := os.Stat(filepath.Join(currentDir, "go.mod")); err == nil {
			return currentDir // Found the root!
		}
		// Move up one directory
		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			// Reached filesystem root, go.mod not found
			return ""
		}
		currentDir = parentDir
	}
}
