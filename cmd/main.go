// cmd/server/main.go  (или просто main.go)
package main

import (
	"context"
	"fmt"
	"log"

	"dailycards/internal/database" // ← sqlc‑generated Queries + Connect()
	"dailycards/internal/server"
	"dailycards/internal/setup"
)

func main() {
	// 1) читаем переменные окружения (.env, docker‑secret и т.д.)
	env := setup.SetupEnv()

	// 2) формируем DSN для Postgres (db — это имя сервиса в docker‑compose)
	dsn := fmt.Sprintf(
		"postgres://%s:%s@db:5432/%s?sslmode=disable",
		env.POSTGRES_USER,
		env.POSTGRES_PASSWORD,
		env.POSTGRES_DB,
	)

	// 3) создаём пул с PreferSimpleProtocol = true (см. database.Connect)
	ctx := context.Background()
	pool, err := database.Connect(ctx, dsn)
	if err != nil {
		log.Fatalf("db connect error: %v", err)
	}
	defer pool.Close()

	// 4) sqlc‑обёртка Queries (умеет работать с *pgxpool.Pool)
	queries := database.New(pool)

	// 5) поднимаем HTTP‑сервер
	srv := server.New(queries, env.SECRET)
	srv.Setup()

	log.Println("⇨ listening on :8080")
	if err := srv.Serve(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}