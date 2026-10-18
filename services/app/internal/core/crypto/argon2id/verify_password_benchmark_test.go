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

	b.ResetTimer()

	for b.Loop() {
		verified, err := hasher.Verify(benchmarkPassword, passwordHash)
		if err != nil {
			b.Fatal(err)
		}
	}
}
