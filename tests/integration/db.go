package integration

import (
	"fmt"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"roadmap.restapi/internal/postgres"
	redisConn "roadmap.restapi/internal/redis"
)

func PostgresConnection(t *testing.T) *sqlx.DB {
	db, err := postgres.Connect("host=localhost port=19876 user=tests password=tests database=tests")
	if err != nil {
		t.Fatal(err.Error())
	}

	return db
}
func CleanTables(t *testing.T, db *sqlx.DB, tables []string) {
	for _, table := range tables {
		db.MustExec(fmt.Sprintf("TRUNCATE %s CASCADE", table))
	}

	db.Close()
}

func RedisConnection(t *testing.T) *redis.Client {
	db, err := redisConn.Connect("127.0.0.1:19875", "", "", 0)
	if err != nil {
		t.Fatal(err.Error())
	}

	return db
}

func RedisClose(t *testing.T, db *redis.Client) {
	db.Close()
}
