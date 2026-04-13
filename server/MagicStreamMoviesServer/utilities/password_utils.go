package utilities

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil{
		return "", err
	}
	return string(encryptedPassword), nil
}

func ComparePassword(hasedPassword string, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hasedPassword), []byte(password))
	if err != nil{
		return err
	}
	return nil
}