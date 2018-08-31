package utils

import (
	"golang.org/x/crypto/bcrypt"
)

type Hash struct{}

func (c *Hash) Generate(s string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	hash := string(hashedBytes[:])
	return hash, nil
}

func (c *Hash) Compare(hash string, s string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(s))
}