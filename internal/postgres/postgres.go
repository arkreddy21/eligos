package postgres

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
)

type DB struct {
	dbpool *pgxpool.Pool
}

func NewDB() *DB {
	dburl, ok := os.LookupEnv("ELIGOSDBURL")
	if !ok {
		log.Fatal("ELIGOSDBURL env variable not set")
	}
	dbpool, err := pgxpool.New(context.Background(), dburl)
	if err != nil {
		log.Fatal("Unable to connect to database: ", err)
	}
	return &DB{
		dbpool: dbpool,
	}
}

func (db *DB) Close() {
	db.dbpool.Close()
}
