package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func GetPasswordHashed(pwd string) string {
	pwdHashed, _ := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)

	return string(pwdHashed)
}

func IsPasswordPassed(pwdHashed, pwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(pwdHashed), []byte(pwd))

	return err == nil
}

func GenerateToken(extra map[string]interface{}) (string, error) {
	claims := jwt.MapClaims{}

	for k, v := range extra {
		claims[k] = v
	}

	// Add exp
	claims["exp"] = jwt.NewNumericDate(time.Now().Add(24 * time.Hour))

	// Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(os.Getenv("SECRET")))
}
