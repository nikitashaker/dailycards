package main

import (
	"context"
	"fmt"
	"log"

	"dailycards/internal/database" 
	"dailycards/internal/server"
	"dailycards/internal/setup"
)

func main() {
	env := setup.SetupEnv()

	dsn := fmt.Sprintf(
		"postgres://%s:%s@db:5432/%s?sslmode=disable",
		env.POSTGRES_USER,
		env.POSTGRES_PASSWORD,
		env.POSTGRES_DB,
	)

	ctx := context.Background()
	pool, err := database.Connect(ctx, dsn)
	if err != nil {
		log.Fatalf("db connect error: %v", err)
	}
	defer pool.Close()

	queries := database.New(pool)

	srv := server.New(queries, env.SECRET)
	srv.Setup()

	log.Println("⇨ listening on :8080")
	if err := srv.Serve(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}