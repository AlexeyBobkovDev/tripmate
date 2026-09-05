package core_crypto_argon2id

import "github.com/alexedwards/argon2id"

func (h *Argon2IDHash) Verify(password, hash string) (bool, error) {
	return argon2id.ComparePasswordAndHash(password, hash)
}
