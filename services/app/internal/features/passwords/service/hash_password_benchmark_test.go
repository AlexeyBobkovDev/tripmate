package passwords_service_test

import (
	"testing"

	"github.com/AlexeyBobkovDev/tripmate/services/app/internal/features/passwords"
	passwords_service "github.com/AlexeyBobkovDev/tripmate/services/app/internal/features/passwords/service"
)

const (
	BENCHMARK_PASSWORD = "BenchmarkPassword_2026!SecureHashTest"
)

func BenchmarkHash(b *testing.B) {
	cfg := passwords.NewConfigMust()
	hasher := passwords_service.NewArgon2IDHasher(
		cfg.Times,
		cfg.Memory,
		cfg.Threads,
		cfg.KeyLen,
	)

	password := []byte(BENCHMARK_PASSWORD)

	b.ResetTimer()

	for b.Loop() {
		_, _, err := hasher.Hash(password)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkVerify(b *testing.B) {
	cfg := passwords.NewConfigMust()
	hasher := passwords_service.NewArgon2IDHasher(
		cfg.Times,
		cfg.Memory,
		cfg.Threads,
		cfg.KeyLen,
	)

	password := []byte(BENCHMARK_PASSWORD)
	passwordHash := []byte{84, 173, 254, 235, 246, 104, 159, 225, 1, 45, 34, 26, 242, 36, 240, 9, 32, 82, 190, 73, 146, 67, 142, 29, 34, 77, 102, 27, 162, 168, 27, 220}
	salt := []byte{154, 9, 108, 87, 240, 107, 240, 110, 27, 213, 192, 160, 164, 173, 215, 132}

	b.ResetTimer()

	for b.Loop() {
		hasher.Verify(password, passwordHash, salt)
	}
}
