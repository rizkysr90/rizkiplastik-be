package main

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rizkysr90/rizkiplastik-be/internal/config"
	"github.com/rizkysr90/rizkiplastik-be/internal/handler"
)

func main() {
	ctx := context.Background()

	// -----------------------------------------------------------------------------------------------------------------
	// LOAD APPLICATION CONFIG FROM ENVIRONMENT VARIABLES
	// -----------------------------------------------------------------------------------------------------------------
	cfg, err := config.LoadFromEnv()
	if err != nil {
		log.Fatalf("main: failed to load and parse config: %s", err)
		return
	}

	// Configure logging to file
	f, _ := createLogFile(cfg.LogPath)
	log.SetOutput(f)

	// -----------------------------------------------------------------------------------------------------------------
	// INFRASTRUCTURE OBJECTS
	// -----------------------------------------------------------------------------------------------------------------
	// Initialize PostgreSQL connection
	dbpool, err := setupDatabase(ctx, cfg.PostgreSQL)
	if err != nil {
		log.Fatalf("main: failed to setup database connection: %s", err)
		return
	}
	defer dbpool.Close()

	// Start the server
	server := handler.NewServer(dbpool)

	log.Printf("Server starting on port %s", cfg.ServerPort)
	if err := server.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// setupDatabase initializes and returns a PostgreSQL connection pool
func setupDatabase(ctx context.Context, pgConfig config.PostgreSQLConfig) (*pgxpool.Pool, error) {
	connString := pgConfig.GetConnectionString()

	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}

	// Set pool configuration
	config.MaxConns = pgConfig.MaxConns
	config.MinConns = pgConfig.MinConns
	config.MaxConnLifetime = pgConfig.MaxConnLifetime
	config.MaxConnIdleTime = pgConfig.MaxConnIdleTime

	// Create the connection pool
	pool, err := pgxpool.ConnectConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	// Test connection
	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	log.Println("Successfully connected to PostgreSQL database")
	return pool, nil
}

// createLogFile creates a log file
func createLogFile(path string) (*os.File, error) {
	if path == "" {
		path = "application.log"
	}
	return os.Create(path)
}
