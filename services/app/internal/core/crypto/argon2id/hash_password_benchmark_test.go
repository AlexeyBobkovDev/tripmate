package core_crypto_argon2id_test

import (
	"testing"

	core_crypto_argon2id "github.com/AlexeyBobkovDev/tripmate/services/app/internal/core/crypto/argon2id"
)

func BenchmarkHash(b *testing.B) {
	hasher := core_crypto_argon2id.NewArgon2IDHash(
		memory,
		times,
		threads,
		saltLen,
		keyLen,
	)

	b.ResetTimer()

	for b.Loop() {
		if _, err := hasher.Hash(BENCHMARK_PASSWORD); err != nil {
			b.Fatal(err)
		}
	}
}
