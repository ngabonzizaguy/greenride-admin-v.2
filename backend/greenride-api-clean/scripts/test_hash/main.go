package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	password := "admin123"
	salt := "f7b1b9f9ea7086d1068c054a459343aa8e3d6ead062878eafdd115fc643b4e0c"
	expectedHash := "$2a$10$1lt6XdTY3Uslh4otXRjz.eocryUsuSCU6lp3lrIlOjao6pk2xaQX."

	// Logic from utils.VerifyPassword
	saltedPassword := password + salt
	preHash := sha256.Sum256([]byte(saltedPassword))
	preHashHex := hex.EncodeToString(preHash[:])

	fmt.Printf("Pre-hash hex: %s\n", preHashHex)

	err := bcrypt.CompareHashAndPassword([]byte(expectedHash), []byte(preHashHex))
	if err == nil {
		fmt.Println("Verification SUCCESS: password matches")
	} else {
		fmt.Printf("Verification FAILED: %v\n", err)
	}
}
