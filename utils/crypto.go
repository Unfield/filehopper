package utils

import "github.com/alexedwards/argon2id"

func HashPassword(input string) (string, error) {
	params := &argon2id.Params{
		Memory:      256 * 1024,
		Iterations:  3,
		Parallelism: 4,
		SaltLength:  16,
		KeyLength:   32,
	}

	return argon2id.CreateHash(input, params)
}

func ComparePassword(plain, hashed string) bool {
	match, err := argon2id.ComparePasswordAndHash(plain, hashed)
	if err != nil {
		return false
	}
	return match
}
