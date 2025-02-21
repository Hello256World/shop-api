package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	pass, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	if err != nil {
		return "", err
	}

	return string(pass), nil
}

func CheckHashPass(newPass, hashPass string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPass), []byte(newPass))
	return err == nil
}
