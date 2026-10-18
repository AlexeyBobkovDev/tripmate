package core_crypto_argon2id_test

import (
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

	for b.Loop() {
		if _, err := hasher.Hash(benchmarkPassword); err != nil {
			b.Fatal(err)
		}
	}
}
