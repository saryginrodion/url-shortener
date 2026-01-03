package integration

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"roadmap.restapi/internal/token"
)

func TestTokenRepo_AddExists(t *testing.T) {
	db := RedisConnection(t)
	defer RedisClose(t, db)

	ctx := context.Background()
	id := uuid.New()
	repo := token.NewRedisWhitelistRepository(db, time.Duration(10) * time.Second)
	err := repo.Add(ctx, id)

	if err != nil {
		t.Errorf("error on add: err: %s", err.Error())
	}

	exists, err := repo.Exists(ctx, id)
	if err != nil {
		t.Errorf("error on exists: err: %s", err.Error())
	}

	if !exists {
		t.Error("added key does not exists")
	}
}

func TestTokenRepo_NotExists(t *testing.T) {
	db := RedisConnection(t)
	defer RedisClose(t, db)

	ctx := context.Background()
	id := uuid.New()
	repo := token.NewRedisWhitelistRepository(db, time.Duration(10) * time.Second)

	exists, err := repo.Exists(ctx, id)
	if err != nil {
		t.Errorf("error on exists: err: %s", err.Error())
	}

	if exists {
		t.Error("key is not added but exists")
	}
}

func TestTokenRepo_DeleteAdded(t *testing.T) {
	db := RedisConnection(t)
	defer RedisClose(t, db)

	ctx := context.Background()
	id := uuid.New()
	repo := token.NewRedisWhitelistRepository(db, time.Duration(10) * time.Second)
	err := repo.Add(ctx, id)

	if err != nil {
		t.Errorf("error on add: err: %s", err.Error())
	}

	err = repo.Remove(ctx, id)
	if err != nil {
		t.Errorf("error on removing: err: %s", err.Error())
	}

	exists, err := repo.Exists(ctx, id)
	if err != nil {
		t.Errorf("error on exists: err: %s", err.Error())
	}

	if exists {
		t.Error("removed key exists")
	}
}
