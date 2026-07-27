package domain

type Password struct {
	UserID  int
	Version int
	Hash    []byte
	Salt    []byte
	Times   uint32
	Memory  uint32
	Threads uint8
	KeyLen  uint32
}

func NewPassword(
	userID int,
	version int,
	hash []byte,
	salt []byte,
	times uint32,
	memory uint32,
	threads uint8,
	keyLen uint32,
) *Password {
	return &Password{
		UserID:  userID,
		Version: version,
		Hash:    hash,
		Salt:    salt,
		Times:   times,
		Memory:  memory,
		Threads: threads,
		KeyLen:  keyLen,
	}
}

func NewPasswordUninitialized(
	hash []byte,
	salt []byte,
	times uint32,
	memory uint32,
	threads uint8,
	keyLen uint32,
) *Password {
	return &Password{
		UserID:  UninitializedID,
		Version: UninitializedVersion,
		Hash:    hash,
		Salt:    salt,
		Times:   times,
		Memory:  memory,
		Threads: threads,
		KeyLen:  keyLen,
	}
}
