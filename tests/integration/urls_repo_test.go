package integration

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"roadmap.restapi/internal/url"
	"roadmap.restapi/internal/user"
)

func createUser(t *testing.T, db *sqlx.DB) uuid.UUID {
	userRepo := user.NewPostgresUserRepository(db)
	users := user.NewUseCases(userRepo, user.NewArgon2IDPasswordHasher())
	user, err := users.NewUser(context.Background(), "test", "test")
	if err != nil {
		t.Fatal(err)
	}

	return user.ID
}

func TestURLRepository_Create_Success(t *testing.T) {
	db := PostgresConnection(t)
	defer CleanTables(t, db, []string{"urls", "users"})
	uid := createUser(t, db)

	repo := url.NewPostgresURLRepository(db)
	ctx := context.Background()

	u := &url.URL{
		ID:       uuid.NewString(),
		AuthorID: uid,
		URL:      "https://example.com",
		Name:     "Example",
	}

	if err := repo.Create(ctx, u); err != nil {
		t.Fatalf("failed to create URL: %v", err)
	}

	if u.ID == "" {
		t.Error("ID should not be empty")
	}
	if u.AuthorID != uid {
		t.Errorf("AuthorID mismatch: got %s, want: %s", u.AuthorID.String(), uid.String())
	}
	if u.URL != "https://example.com" {
		t.Errorf("url mismatch: got %s", u.URL)
	}
	if u.Name != "Example" {
		t.Errorf("name mismatch: got %s", u.Name)
	}
	if time.Since(u.CreatedAt) > time.Second {
		t.Errorf("CreatedAt is too old: %v", u.CreatedAt)
	}
}

func TestURLRepository_Create_DuplicateID(t *testing.T) {
	db := PostgresConnection(t)
	defer CleanTables(t, db, []string{"urls", "users"})
	uid := createUser(t, db)

	repo := url.NewPostgresURLRepository(db)
	ctx := context.Background()

	id := uuid.NewString()

	u1 := &url.URL{
		ID:       id,
		AuthorID: uid,
		URL:      "https://one.test",
		Name:     "One",
	}
	if err := repo.Create(ctx, u1); err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	u2 := &url.URL{
		ID:       id,
		AuthorID: uid,
		URL:      "https://two.test",
		Name:     "Two",
	}
	err := repo.Create(ctx, u2)
	if err == nil {
		t.Fatal("expected error on duplicate id, got nil")
	}
	if err != url.ErrURLAlreadyExists {
		t.Errorf("wrong error: got %v, want %v", err, url.ErrURLAlreadyExists)
	}
}

func TestURLRepository_ByID_Success(t *testing.T) {
	db := PostgresConnection(t)
	defer CleanTables(t, db, []string{"urls", "users"})
	uid := createUser(t, db)

	repo := url.NewPostgresURLRepository(db)
	ctx := context.Background()

	u := &url.URL{
		ID:       uuid.NewString(),
		AuthorID: uid,
		URL:      "https://byid.test",
		Name:     "ByID",
	}

	if err := repo.Create(ctx, u); err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	found, err := repo.ByID(ctx, u.ID)
	if err != nil {
		t.Fatalf("failed to get by id: %v", err)
	}

	if found.ID != u.ID {
		t.Errorf("ID mismatch: got %s, want %s", found.ID, u.ID)
	}
	if found.URL != u.URL {
		t.Errorf("url mismatch: got %s, want %s", found.URL, u.URL)
	}
	if found.Name != u.Name {
		t.Errorf("name mismatch: got %s, want %s", found.Name, u.Name)
	}
}

func TestURLRepository_ByID_NotFound(t *testing.T) {
	db := PostgresConnection(t)
	defer CleanTables(t, db, []string{"urls", "users"})

	repo := url.NewPostgresURLRepository(db)
	ctx := context.Background()

	_, err := repo.ByID(ctx, uuid.NewString())
	if err == nil {
		t.Fatal("expected ErrURLNotFound, got nil")
	}
	if err != url.ErrURLNotFound {
		t.Errorf("wrong error: got %v, want %v", err, url.ErrURLNotFound)
	}
}

func TestURLRepository_Update_Success(t *testing.T) {
	db := PostgresConnection(t)
	defer CleanTables(t, db, []string{"urls", "users"})
	uid := createUser(t, db)

	repo := url.NewPostgresURLRepository(db)
	ctx := context.Background()

	u := &url.URL{
		ID:       uuid.NewString(),
		AuthorID: uid,
		URL:      "https://old.test",
		Name:     "Old",
	}

	if err := repo.Create(ctx, u); err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	u.URL = "https://new.test"
	u.Name = "New"

	if err := repo.Update(ctx, u); err != nil {
		t.Fatalf("failed to update: %v", err)
	}

	updated, err := repo.ByID(ctx, u.ID)
	if err != nil {
		t.Fatalf("failed to reload: %v", err)
	}

	if updated.URL != "https://new.test" {
		t.Errorf("url not updated: got %s", updated.URL)
	}
	if updated.Name != "New" {
		t.Errorf("name not updated: got %s", updated.Name)
	}
}

func TestURLRepository_Update_NotFound(t *testing.T) {
	db := PostgresConnection(t)
	defer CleanTables(t, db, []string{"urls", "users"})
	uid := createUser(t, db)

	repo := url.NewPostgresURLRepository(db)
	ctx := context.Background()

	u := &url.URL{
		ID:       uuid.NewString(),
		AuthorID: uid,
		URL:      "https://ghost.test",
		Name:     "Ghost",
	}

	err := repo.Update(ctx, u)
	if err == nil {
		t.Fatal("expected error on update non-existent url, got nil")
	}
	if err != url.ErrURLNotFound {
		t.Errorf("wrong error: got %v, want %v", err, url.ErrURLNotFound)
	}
}

func TestURLRepository_Delete_Success(t *testing.T) {
	db := PostgresConnection(t)
	defer CleanTables(t, db, []string{"urls", "users"})
	uid := createUser(t, db)

	repo := url.NewPostgresURLRepository(db)
	ctx := context.Background()

	u := &url.URL{
		ID:       uuid.NewString(),
		AuthorID: uid,
		URL:      "https://delete.test",
		Name:     "Delete",
	}

	if err := repo.Create(ctx, u); err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	if err := repo.Delete(ctx, u.ID); err != nil {
		t.Fatalf("failed to delete: %v", err)
	}

	_, err := repo.ByID(ctx, u.ID)
	if err != url.ErrURLNotFound {
		t.Errorf("expected ErrURLNotFound after delete, got %v", err)
	}
}

func TestURLRepository_Delete_NotFound(t *testing.T) {
	db := PostgresConnection(t)
	defer CleanTables(t, db, []string{"urls", "users"})

	repo := url.NewPostgresURLRepository(db)
	ctx := context.Background()

	err := repo.Delete(ctx, uuid.NewString())
	if err != url.ErrURLNotFound {
		t.Errorf("expected ErrURLNotFound on delete non-existent, got %v", err)
	}
}

func TestURLRepository_ByUser_Success(t *testing.T) {
	db := PostgresConnection(t)
	defer CleanTables(t, db, []string{"urls", "users"})
	uid := createUser(t, db)

	repo := url.NewPostgresURLRepository(db)
	ctx := context.Background()

	u1 := &url.URL{
		ID:       uuid.NewString(),
		AuthorID: uid,
		URL:      "https://one.test",
		Name:     "One",
	}
	u2 := &url.URL{
		ID:       uuid.NewString(),
		AuthorID: uid,
		URL:      "https://two.test",
		Name:     "Two",
	}

	if err := repo.Create(ctx, u1); err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	if err := repo.Create(ctx, u2); err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	list, err := repo.ByUser(ctx, uid)
	if err != nil {
		t.Fatalf("failed to get by user: %v", err)
	}

	if len(list) != 2 {
		t.Fatalf("expected 2 urls, got %d", len(list))
	}
}

