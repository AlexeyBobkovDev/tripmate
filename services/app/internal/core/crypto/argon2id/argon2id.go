package core_crypto_argon2id

import "github.com/alexedwards/argon2id"

type Argon2IDHash struct {
	*argon2id.Params
}

func NewArgon2IDHash(
	Memory uint32,
	Iterations uint32,
	Parallelism uint8,
	SaltLength uint32,
	KeyLength uint32,
) *Argon2IDHash {
	return &Argon2IDHash{
		&argon2id.Params{
			Memory:      Memory,
			Iterations:  Iterations,
			Parallelism: Parallelism,
			SaltLength:  SaltLength,
			KeyLength:   KeyLength,
		},
	}
}
