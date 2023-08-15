package util

import (
	"golang.org/x/crypto/bcrypt"
	"os"
	"strconv"
)

var minCryptoCost = 16

func HashPassword(password string) ([]byte, error) {

	// Hash the password
	return bcrypt.GenerateFromPassword([]byte(password), getCryptoCost())
}

func getCryptoCost() int {
	// Get env settings
	envCost := os.Getenv("CRYPT_COST")

	envIntCost, err := strconv.ParseInt(envCost, 10, 8)
	if err != nil || int(envIntCost) < minCryptoCost {
		return minCryptoCost
	}
	return int(envIntCost)
}
