package crypt

type CryptStrategy interface {
	Hash(password string) (string, error)
	Compare(hashedPassword string, password string) bool
}
