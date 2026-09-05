package core_crypto_argon2id

import (
	"github.com/alexedwards/argon2id"
)

func (h *Argon2IDHash) Hash(password string) (string, error) {
	params := &argon2id.Params{
		Memory:      h.Memory,
		Iterations:  h.Iterations,
		Parallelism: h.Parallelism,
		SaltLength:  h.SaltLength,
		KeyLength:   h.KeyLength,
	}
	return argon2id.CreateHash(password, params)
}
