package helpers

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"errors"
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
	initialVector = "1010101010101010"
	key           = []byte{0x31, 0x31, 0x31, 0x31, 0x31, 0x31, 0x31, 0x31, 0x31, 0x31, 0x31, 0x31, 0x31, 0x31, 0x31, 0x31}
)

func Encrypt(payload string) (string, error) {

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	ecb := cipher.NewCBCEncrypter(block, []byte(initialVector))
	content := []byte(payload)
	content = PKCS5Padding(content, block.BlockSize())
	crypted := make([]byte, len(content))
	ecb.CryptBlocks(crypted, content)
	encryptedString := hex.EncodeToString([]byte(crypted))
	return encryptedString, nil
}

func Decrypt(payload string) (string, error) {

	if payload == "" {
		return "", errors.New("empty string error")
	}

	crypt, err := hex.DecodeString(payload)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	ecb := cipher.NewCBCDecrypter(block, []byte(initialVector))
	decrypted := make([]byte, len(crypt))
	ecb.CryptBlocks(decrypted, crypt)
	decrypted = PKCS5Trimming(decrypted)
	return string(decrypted), nil
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5Trimming(encrypt []byte) []byte {
	padding := encrypt[len(encrypt)-1]
	return encrypt[:len(encrypt)-int(padding)]
}
