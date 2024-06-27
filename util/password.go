package util

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// it should return the bcrypt hash of the password

func HashPassword(password string) (string , error) {
	hashedPassword , err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "" , fmt.Errorf("error in hashing password %w", err)
	}
	return string(hashedPassword) , nil
}

// check the provided password is correct or not

func CheckPassword(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}