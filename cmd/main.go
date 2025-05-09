package main

import (
    "context"
    "fmt"
    "log"

    "github.com/jackc/pgx/v5"

    "dailycards/internal/database"
    "dailycards/internal/server"
    "dailycards/internal/setup"
)

func main() {
    // 1) Загружаем переменные окружения (POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_DB, SECRET)
    env := setup.SetupEnv()

    // 2) Подключаемся к Postgres
    ctx := context.Background()
    dsn := fmt.Sprintf(
        "postgresql://%v:%v@db:5432/%v?sslmode=disable",
        env.POSTGRES_USER,
        env.POSTGRES_PASSWORD,
        env.POSTGRES_DB,
    )
    conn, err := pgx.Connect(ctx, dsn)
    if err != nil {
        log.Fatalf("failed to connect to database: %v", err)
    }
    defer conn.Close(ctx)

    // 3) Инициализируем sqlc-обёртку
    db := database.New(conn)

    // 4) Создаём сервер (единственный echo.Echo + session-middleware внутри)
    srv := server.New(db, env.SECRET)

    // 5) Регистрируем все маршруты и middleware
    srv.Setup()

    // 6) Запускаем на :8080
    if err := srv.Serve(); err != nil {
        log.Fatalf("server error: %v", err)
    }
}

