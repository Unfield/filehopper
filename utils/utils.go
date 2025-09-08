package utils

import gonanoid "github.com/matoous/go-nanoid/v2"

func GenerateNanoid() (string, error) {
	return gonanoid.Generate("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-_", 30)
}
