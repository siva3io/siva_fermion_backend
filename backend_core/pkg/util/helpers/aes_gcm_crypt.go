package helpers

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"os"
)

var (
	AES_SECRET_KEY = os.Getenv("AES_SECRET_KEY")
)

func AESGCMEncrypt(text string) (string, error) {
	dStr, _ := hex.DecodeString(AES_SECRET_KEY)
	plaintext := []byte(text)

	hasher := sha512.New()
	hasher.Write(dStr)
	out := hex.EncodeToString(hasher.Sum(nil))
	newKey, _ := hex.DecodeString(out[:64])
	nonce, _ := hex.DecodeString(out[64:(64 + 24)])

	block, err := aes.NewCipher(newKey)
	if err != nil {
		return "", err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	cipherText := aesgcm.Seal(nil, nonce, plaintext, nil)
	return fmt.Sprintf("%x", cipherText), nil
}

func AESGCMDecrypt(data string) (string, error) {
	dStr, _ := hex.DecodeString(AES_SECRET_KEY)

	hasher := sha512.New()
	hasher.Write(dStr)
	out := hex.EncodeToString(hasher.Sum(nil))
	newKey, _ := hex.DecodeString(out[:64])
	nonce, _ := hex.DecodeString(out[64:(64 + 24)])

	block, err := aes.NewCipher(newKey)
	if err != nil {
		return "", err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	cipherText, _ := hex.DecodeString(data)
	output, _ := aesgcm.Open(nil, nonce, cipherText, nil)

	return string(output), nil
}
