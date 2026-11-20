package core_crypto_argon2id_test

import (
	"testing"

	core_crypto_argon2id "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/crypto/argon2id"
)

func BenchmarkVerify(b *testing.B) {
	hasher := core_crypto_argon2id.NewArgon2IDHash(
		memoryKiB,
		iterations,
		parallelism,
		saltLen,
		keyLen,
	)

	passwordHash, err := hasher.Hash(benchmarkPassword)
	if err != nil {
		b.Fatal(err)
	}

	b.Run("correct password", func(b *testing.B) {
		for b.Loop() {
			verified, err := hasher.Verify(benchmarkPassword, passwordHash)
			if err != nil {
				b.Fatal(err)
			}
			if !verified {
				b.Fatal("expected password verification to succeed")
			}
		}
	})

	b.Run("incorrect password", func(b *testing.B) {
		for b.Loop() {
			verified, err := hasher.Verify("WrongPassword", passwordHash)
			if err != nil {
				b.Fatal(err)
			}
			if verified {
				b.Fatal("expected password verification to fail")
			}
		}
	})

	b.Run("malformed hash", func(b *testing.B) {
		for b.Loop() {
			_, err := hasher.Verify(benchmarkPassword, "not-a-valid-argon2id-hash")
			if err == nil {
				b.Fatal("expected an error for a malformed hash")
			}
		}
	})

	b.Run("empty password", func(b *testing.B) {
		for b.Loop() {
			_, err := hasher.Verify("", passwordHash)
			if err == nil {
				b.Fatal("expected an error for empty password")
			}
		}
	})

	b.Run("empty hash", func(b *testing.B) {
		for b.Loop() {
			_, err := hasher.Verify(benchmarkPassword, "")
			if err == nil {
				b.Fatal("expected an error for empty hash")
			}
		}
	})
}
