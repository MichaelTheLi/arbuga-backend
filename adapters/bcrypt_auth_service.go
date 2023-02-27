package adapters

import (
	"golang.org/x/crypto/bcrypt"
)

type BcryptAuthService struct {
}

func (auth BcryptAuthService) HashFromPassword(password string) (string, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)

	return string(hashedPass), err
}

func (auth BcryptAuthService) IsHashValid(expectedHash string, password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(expectedHash), []byte(password))
	if err != nil {
		return false, err
	}

	return true, err
}
