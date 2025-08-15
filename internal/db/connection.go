package db

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var dbInstance *pgx.Conn
var mu sync.Mutex

func Connect() (*pgx.Conn, error) {
	if dbInstance != nil {
		log.Println("⚠️ Reusing existing database connection")
		return dbInstance, nil
	}

	user := "user"
	pass := "password"
	host := "localhost"
	port := "5432"
	name := "test"

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user, pass, host, port, name,
	)

	ctx := context.Background()
	db, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	if err := db.Ping(ctx); err != nil {
		return nil, fmt.Errorf("error pinging database: %w", err)
	}

	log.Println("✅ Connected to PostgreSQL")
	if dbInstance == nil {
		mu.Lock()
		defer mu.Unlock()
		dbInstance = db
	} else {
		log.Println("⚠️ Reusing existing database connection")
		db.Close(ctx)
		return dbInstance, nil
	}

	return dbInstance, nil
}
