package crypt

import "golang.org/x/crypto/bcrypt"

type BcryptStrategy struct {
}

func NewBcryptStrategy() CryptStrategy {
	return &BcryptStrategy{}
}

func (b *BcryptStrategy) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func (b *BcryptStrategy) Compare(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
