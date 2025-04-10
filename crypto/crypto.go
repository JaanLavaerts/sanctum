package crypto

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func GenerateHash(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

func VerifyMasterPassword(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}

func GenerateToken() (raw string, hashed string) {
    // first generate 32 random bytes, then encode them using base64
    // raw gets stored in the cookie, while the hashed version gets stored in the DB
    b := make([]byte, 32)
    _, err := rand.Read(b)
    if err != nil {
        log.Fatal(err)
    }

    raw = base64.RawURLEncoding.EncodeToString(b)

    h := sha256.New()
    h.Write(b)
    hashed = hex.EncodeToString(h.Sum(nil))

    return raw, hashed
}

func VerifyToken(token string) bool {
    // TODO
    // 1. decode the token using base64: string => bytes
    // 2. hash the raw bytes: bytes => string
    // 3. compare this hash to the hash in DB
    return false
}
