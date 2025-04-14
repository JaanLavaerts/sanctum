package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"log"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/scrypt"
)

func GenerateHash(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

func VerifyMasterPassword(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}

func DeriveKey(password string, salt []byte) ([]byte, error) {
    key, err := scrypt.Key([]byte(password), salt, 1<<15, 8, 1, 32)
    return key, err
}

func GenerateSalt() ([]byte, error) {
    salt := make([]byte, 16)
    _, err := rand.Read(salt)
    return salt, err
}

func EncryptEntryPassword(plainPassword string, derrivedKey []byte) ([]byte, []byte, error) {
    nonce := make([]byte, 12)
    _, err := rand.Read(nonce)
    if err != nil {
        return nil, nil, err
    }

    block, err := aes.NewCipher(derrivedKey)
    if err != nil {
        return nil, nil, err
    }

    aesGCM, err := cipher.NewGCM(block)
    if err != nil {
        return nil, nil, err
    }

    encryptedPassword := aesGCM.Seal(nil, nonce, []byte(plainPassword), nil)

    return encryptedPassword, nonce, nil
}

func DecryptPassword(encryptedPassword string, derivedKey []byte, nonce string) (string, error) {
    nonceString, _ := base64.StdEncoding.DecodeString(nonce)
    passwordString, _ := base64.StdEncoding.DecodeString(encryptedPassword)

    block, _ := aes.NewCipher(derivedKey)
    aesGCM, _ := cipher.NewGCM(block)

    plainPassword, err := aesGCM.Open(nil, nonceString, passwordString, nil)
    if err != nil {
        return "", err 
    }
    return string(plainPassword), nil
}

func GenerateAuthToken() (raw string, hashed string) {
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

func VerifyAuthToken(input_token string, db_token string) bool {
    b, err := base64.RawURLEncoding.DecodeString(input_token) 
    if err != nil {
        log.Fatal(err)
    }

    h := sha256.New()
    h.Write(b)
    hashed_bytes := hex.EncodeToString(h.Sum(nil)) 

    return hashed_bytes == db_token 
}