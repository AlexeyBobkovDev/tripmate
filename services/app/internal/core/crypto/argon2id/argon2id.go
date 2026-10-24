package core_crypto_argon2id

import "github.com/alexedwards/argon2id"

type Argon2IDHash struct {
	*argon2id.Params
}

func NewArgon2IDHash(
	memory uint32,
	iterations uint32,
	parallelism uint8,
	saltLength uint32,
	keyLength uint32,
) *Argon2IDHash {
	return &Argon2IDHash{
		&argon2id.Params{
			Memory:      memory,
			Iterations:  iterations,
			Parallelism: parallelism,
			SaltLength:  saltLength,
			KeyLength:   keyLength,
		},
	}
}
