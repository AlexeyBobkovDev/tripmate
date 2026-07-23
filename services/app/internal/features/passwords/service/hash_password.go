package passwords_service

import (
	"crypto/rand"
	"fmt"
	"slices"

	"golang.org/x/crypto/argon2"
)

func (h *Argon2IDHasher) Hash(password []byte) ([]byte, []byte, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return []byte{}, []byte{}, fmt.Errorf("hash password: %w", err)
	}

	hash := hashPassword(
		[]byte(password),
		salt,
		h.Times,
		h.Memory,
		h.Threads,
		h.KeyLen,
	)

	return hash, salt, nil
}

func (h *Argon2IDHasher) Verify(password, encryptedPassword []byte, salt []byte) bool {
	return slices.Equal(hashPassword(password, salt, h.Times, h.Memory, h.Threads, h.KeyLen), encryptedPassword)
}

func hashPassword(password, salt []byte, time, memory uint32, threads uint8, keyLen uint32) []byte {
	return argon2.IDKey(
		[]byte(password),
		salt,
		time,
		memory,
		threads,
		keyLen,
	)
}
