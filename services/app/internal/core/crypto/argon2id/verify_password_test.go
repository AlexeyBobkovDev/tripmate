package core_crypto_argon2id_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	core_crypto_argon2id "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/crypto/argon2id"
)

func TestVerify(t *testing.T) {
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

	t.Run("correct password", func(t *testing.T) {
		verified, err := hasher.Verify(password, passwordHash)
		if err != nil {
			t.Fatal(err)
		}
		if !verified {
			t.Fatal("expected password verification to succeed")
		}
	})

	t.Run("incorrect password", func(t *testing.T) {
		verified, err := hasher.Verify("WrongPassword", passwordHash)
		if err != nil {
			t.Fatal(err)
		}
		if verified {
			t.Fatal("expected password verification to fail")
		}
	})

	t.Run("malformed hash", func(t *testing.T) {
		_, err := hasher.Verify(password, "not-a-valid-argon2id-hash")
		if err == nil {
			t.Fatal("expected an error for a malformed hash")
		}
	})

	t.Run("empty password", func(t *testing.T) {
		// Note: argon2id library doesn't reject empty passwords
		valid, err := hasher.Verify("", passwordHash)
		assert.False(t, valid)
		assert.NoError(t, err)
	})

	t.Run("empty hash", func(t *testing.T) {
		_, err := hasher.Verify(password, "")
		if err == nil {
			t.Fatal("expected an error for empty hash")
		}
	})

	t.Run("case sensitivity", func(t *testing.T) {
		upperPassword := strings.ToUpper(password)
		verified, err := hasher.Verify(upperPassword, passwordHash)
		if err != nil {
			t.Fatal(err)
		}
		if verified {
			t.Fatal("expected password verification to fail for case difference")
		}
	})

	t.Run("similar password", func(t *testing.T) {
		similarPassword := strings.Replace(password, "!", "?", 1)
		verified, err := hasher.Verify(similarPassword, passwordHash)
		if err != nil {
			t.Fatal(err)
		}
		if verified {
			t.Fatal("expected password verification to fail for similar password")
		}
	})

	t.Run("very long password", func(t *testing.T) {
		longPassword := strings.Repeat("a", 1000)
		verified, err := hasher.Verify(longPassword, passwordHash)
		if err != nil {
			t.Fatal(err)
		}
		if verified {
			t.Fatal("expected password verification to fail for very long password")
		}
	})

	t.Run("unicode password", func(t *testing.T) {
		unicodePassword := "密码123пароль"
		verified, err := hasher.Verify(unicodePassword, passwordHash)
		if err != nil {
			t.Fatal(err)
		}
		if verified {
			t.Fatal("expected password verification to fail for unicode password")
		}
	})
}
