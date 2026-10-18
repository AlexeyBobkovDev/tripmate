package core_crypto_argon2id_test

import (
	"testing"

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
}
