package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, "postgresql://root:secret@localhost:5432/todo_app?sslmode=disable")
	if err != nil {
		log.Fatal("pool creation failed: ", err)
	}
	defer pool.Close()

	// pgxpool.New() is lazy it only validates connection string but doesn't actually dial the db thus check ping
	err = pool.Ping(ctx)
	if err != nil {
		log.Fatal("connection failed: ", err)
	}

	log.Print("server connected\n")
}
