package main

import (
	"context"
	"log"
	"os"

	"github.com/MidNight91119/todo-app/api"
	db "github.com/MidNight91119/todo-app/db/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgresql://root:secret@localhost:5432/todo_app?sslmode=disable"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatal("pool creation failed: ", err)
	}
	defer pool.Close()

	// pgxpool.New() is lazy it only validates connection string but doesn't actually dial the db thus check ping
	err = pool.Ping(ctx)
	if err != nil {
		log.Fatal("connection failed: ", err)
	}

	log.Print("pool connected\n")

	queries := db.New(pool)
	server := api.NewServer(queries)

	if err := server.Start(":" + port); err != nil {
		log.Fatal("cannot start server: ", err)
	}
}
