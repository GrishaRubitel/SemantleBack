package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/url"
	"strings"
)

func createKey() string {

	var keyCont []byte = []byte("this's secret key.enough 32 bits")
	if _, err := rand.Read(keyCont); err != nil {
		panic(err.Error())
	}
	return hex.EncodeToString(keyCont)
}

func encrypt(stringToEncrypt string) (encryptedString string) {

	key, _ := hex.DecodeString(siphKey)
	plaintext := []byte(stringToEncrypt)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return base64.URLEncoding.EncodeToString(ciphertext)
}

func decrypt(stringToDecrypt string) string {

	key, _ := hex.DecodeString(siphKey)
	ciphertext, _ := base64.URLEncoding.DecodeString(stringToDecrypt)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	if len(ciphertext) < aes.BlockSize {
		panic("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	stream.XORKeyStream(ciphertext, ciphertext)

	return fmt.Sprintf("%s", ciphertext)
}

func urlDecoder(target string, flag int) string {
	deTarget, err := url.QueryUnescape(target)
	if err != nil {
		log.Fatalln(err)
		return err.Error()
	}
	if flag == 0 {
		return decrypt(deTarget)
	}
	return strings.ToLower(deTarget)
}
