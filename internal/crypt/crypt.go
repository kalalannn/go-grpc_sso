package crypt

import "golang.org/x/crypto/bcrypt"

func HashedPassword(plainPassword string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	return string(hashed), err
}

func CMPHashedPlainPassword(hashedPassword, plainPassword string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword)) == nil
}
