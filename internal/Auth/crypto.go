package Auth

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"golang.org/x/crypto/pbkdf2"
)

func PasswordPBKDF2(plainPassword []byte, sault []byte) []byte {
	return pbkdf2.Key(plainPassword, sault, 4096, 32, sha1.New)
}

func GenSault(login string) []byte {
	return md5.New().Sum([]byte(login))
}

func CheckPassword(toVerify []byte, gage []byte, sault []byte) bool{
	hashVerify := PasswordPBKDF2(toVerify, sault)
	return bytes.Equal(gage, hashVerify)
}