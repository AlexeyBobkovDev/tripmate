package core_crypto_argon2id_test

import (
	"strings"
	"testing"

	core_crypto_argon2id "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/crypto/argon2id"
)

func BenchmarkHash(b *testing.B) {
	hasher := core_crypto_argon2id.NewArgon2IDHash(
		memoryKiB,
		iterations,
		parallelism,
		saltLen,
		keyLen,
	)

	b.Run("standard password", func(b *testing.B) {
		for b.Loop() {
			if _, err := hasher.Hash(benchmarkPassword); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("empty password", func(b *testing.B) {
		for b.Loop() {
			if _, err := hasher.Hash(""); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("short password", func(b *testing.B) {
		shortPassword := "short"
		for b.Loop() {
			if _, err := hasher.Hash(shortPassword); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("long password", func(b *testing.B) {
		longPassword := strings.Repeat("a", 1000)
		for b.Loop() {
			if _, err := hasher.Hash(longPassword); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("unicode password", func(b *testing.B) {
		unicodePassword := "密码123пароль"
		for b.Loop() {
			if _, err := hasher.Hash(unicodePassword); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("special characters", func(b *testing.B) {
		specialPassword := "P@ssw0rd!@#$%^&*()_+-={}[]|:;'<>,.?/"
		for b.Loop() {
			if _, err := hasher.Hash(specialPassword); err != nil {
				b.Fatal(err)
			}
		}
	})
}
