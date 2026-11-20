package core_crypto_argon2id_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	core_crypto_argon2id "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/crypto/argon2id"
)

func TestHash(t *testing.T) {
	hasher := core_crypto_argon2id.NewArgon2IDHash(
		memoryKiB,
		iterations,
		parallelism,
		saltLen,
		keyLen,
	)

	passwordHash, err := hasher.Hash(password)
	if err != nil {
		t.Fatal(err)
	}
	if passwordHash == "" {
		t.Fatal("expected a non-empty password hash")
	}
}

func TestHash_EmptyPassword(t *testing.T) {
	hasher := core_crypto_argon2id.NewArgon2IDHash(
		memoryKiB,
		iterations,
		parallelism,
		saltLen,
		keyLen,
	)

	// Note: argon2id library doesn't reject empty passwords
	_, err := hasher.Hash("")
	assert.NoError(t, err)
}

func TestHash_VeryLongPassword(t *testing.T) {
	hasher := core_crypto_argon2id.NewArgon2IDHash(
		memoryKiB,
		iterations,
		parallelism,
		saltLen,
		keyLen,
	)

	longPassword := strings.Repeat("a", 10000)
	passwordHash, err := hasher.Hash(longPassword)
	if err != nil {
		t.Fatal(err)
	}
	if passwordHash == "" {
		t.Fatal("expected a non-empty password hash for long password")
	}
}

func TestHash_SpecialCharacters(t *testing.T) {
	hasher := core_crypto_argon2id.NewArgon2IDHash(
		memoryKiB,
		iterations,
		parallelism,
		saltLen,
		keyLen,
	)

	specialPassword := "P@ssw0rd!@#$%^&*()_+-={}[]|:;'<>,.?/"
	passwordHash, err := hasher.Hash(specialPassword)
	if err != nil {
		t.Fatal(err)
	}
	if passwordHash == "" {
		t.Fatal("expected a non-empty password hash for special characters")
	}
}

func TestHash_SamePasswordDifferentHash(t *testing.T) {
	hasher := core_crypto_argon2id.NewArgon2IDHash(
		memoryKiB,
		iterations,
		parallelism,
		saltLen,
		keyLen,
	)

	hash1, err := hasher.Hash(password)
	if err != nil {
		t.Fatal(err)
	}

	hash2, err := hasher.Hash(password)
	if err != nil {
		t.Fatal(err)
	}

	if hash1 == hash2 {
		t.Fatal("expected different hashes for the same password (different salts)")
	}
}

func TestHash_UnicodePassword(t *testing.T) {
	hasher := core_crypto_argon2id.NewArgon2IDHash(
		memoryKiB,
		iterations,
		parallelism,
		saltLen,
		keyLen,
	)

	unicodePassword := "密码123пароль"
	passwordHash, err := hasher.Hash(unicodePassword)
	if err != nil {
		t.Fatal(err)
	}
	if passwordHash == "" {
		t.Fatal("expected a non-empty password hash for unicode password")
	}
}
