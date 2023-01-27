package helpers

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"os"
)
/*
Copyright (C) 2022 Eunimart Omnichannel Pvt Ltd. (www.eunimart.com)
All rights reserved.
This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Lesser General Public License v3.0 as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.
This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Lesser General Public License v3.0 for more details.
You should have received a copy of the GNU Lesser General Public License v3.0
along with this program.  If not, see <https://www.gnu.org/licenses/lgpl-3.0.html/>.
*/
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
