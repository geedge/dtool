package utils

import (
	"bytes"
	"crypto/des"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
)

func pkcs5Padding(plaintext []byte, blocksize int) []byte {
	padding := blocksize - len(plaintext)%blocksize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(plaintext, padtext...)
}

func pkcs5Unpadding(plaintext []byte) []byte {
	padding := plaintext[len(plaintext)-1]
	return plaintext[:len(plaintext)-int(padding)]
}

func encryptDES(plaintext, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blocksize := block.BlockSize()
	plaintext = pkcs5Padding(plaintext, blocksize)

	ciphertext := make([]byte, len(plaintext))
	for i := 0; i < len(plaintext); i += blocksize {
		block.Encrypt(ciphertext[i:], plaintext[i:i+blocksize])
	}

	return ciphertext, nil
}

func decryptDES(ciphertext, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blocksize := block.BlockSize()
	decryptedtext := make([]byte, len(ciphertext))
	for i := 0; i < len(ciphertext); i += blocksize {
		block.Decrypt(decryptedtext[i:], ciphertext[i:i+blocksize])
	}
	decryptedtext = pkcs5Unpadding(decryptedtext)
	return decryptedtext, nil
}

func SendTencentHttpdnsQuery() {
	client := &http.Client{}
	domain := []byte("echo.echodns.xyz")
	key := []byte("046Ju3Cw")
	encrypted_bytes, err := encryptDES(domain, key)
	if err != nil {
		return
	}
	encrypted_domain := hex.EncodeToString(encrypted_bytes)
	url := fmt.Sprintf("http://119.29.29.98/d?dn=%s&id=61188", encrypted_domain)
	fmt.Println(url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("create new request failed. Error: %s\n", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("request went wrong. Error: %s\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("read content failed. Error: %s\n", err)
		return
	}
	encrypted_response, _ := hex.DecodeString(string(body))
	decrypted_response, err := decryptDES(encrypted_response, key)
	if err != nil {
		return
	}
	fmt.Println(string(decrypted_response))
}

func SendAlicloudHttpdnsQurey() {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://203.107.1.33/149702/d?host=echo.echodns.xyz", nil)
	if err != nil {
		fmt.Printf("create new request failed. Error: %s\n", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("request went wrong. Error: %s\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("read content failed. Error: %s\n", err)
		return
	}

	fmt.Println(string(body))
}
