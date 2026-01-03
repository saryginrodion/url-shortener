package user

import (
	"fmt"

	"github.com/alexedwards/argon2id"
)


type Argon2IDPasswordHasher struct { }

func NewArgon2IDPasswordHasher() *Argon2IDPasswordHasher {
	return &Argon2IDPasswordHasher{}
}

func (h *Argon2IDPasswordHasher) Hash(password string) string {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		panic(fmt.Sprintf("failed to hash password: err: %s", err.Error()))
	}

	return hash
}

func (h *Argon2IDPasswordHasher) Check(password string, hash string) bool {
	matched, _, err := argon2id.CheckHash(password, hash)
	if err != nil {
		panic(fmt.Sprintf("failed to check password hash: err: %s", err.Error()))
	}

	return matched
}
