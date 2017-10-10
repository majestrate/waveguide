package util

import (
	"encoding/base64"
	"golang.org/x/crypto/bcrypt"
)

// CheckLogin checks password against login credential
// returns true if the password matches otherwise false
func CheckLogin(login, password string) bool {
	hashed, _ := base64.StdEncoding.DecodeString(login)
	return bcrypt.CompareHashAndPassword(hashed, []byte(password)) == nil
}

func NewPassword(password string) string {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return base64.StdEncoding.EncodeToString(hashed)
}
