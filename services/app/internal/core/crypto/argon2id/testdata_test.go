package core_crypto_argon2id_test

const (
	memoryKiB   uint32 = 64 * 1024
	iterations  uint32 = 3
	parallelism uint8  = 2
	saltLen     uint32 = 16
	keyLen      uint32 = 32

	benchmarkPassword = "BenchmarkPassword_2026!SecureHashTest"
	password          = "Password_2026!SecureHashTest"
)
