package unit

import (
	"strings"
	"testing"

	"roadmap.restapi/internal/user"
)

func TestArgon2IDHasher_CreateAndCheckMatching(t *testing.T) {
	hasher := user.NewArgon2IDPasswordHasher()

	pwd := "test-super-secret-pwd"
	hash := hasher.Hash(pwd)

	if len(strings.TrimSpace(hash)) == 0 {
		t.Fatal("hash is empty")
	}

	if !hasher.Check(pwd, hash) {
		t.Fatal("check failed")
	}
}

func TestArgon2IDHasher_CheckUnmatching(t *testing.T) {
	hasher := user.NewArgon2IDPasswordHasher()

	pwd := "test-super-secret-pwd"
	hash := hasher.Hash(pwd)

	if len(strings.TrimSpace(hash)) == 0 {
		t.Fatal("hash is empty")
	}

	if hasher.Check("incorrect-password", hash) {
		t.Fatal("different password matched")
	}
}
