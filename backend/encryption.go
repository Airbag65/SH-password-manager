package main

import (
	"crypto/sha256"
	"fmt"
)

func extractSalt(pwd string) (string, string) {
	saltBeginning := ""
	for i := len(pwd) - 2; i < len(pwd); i++ {
		saltBeginning = fmt.Sprintf("%s%c", saltBeginning, pwd[i])
	}

	saltEnd := ""
	for _, c := range pwd[:5] {
		saltEnd = fmt.Sprintf("%s%c", saltEnd, c)
	}
	return saltBeginning, saltEnd
}

func encryptPassword(origPwd string) string {
	encPassword := origPwd
	saltBeginning, saltEnd := extractSalt(encPassword)

	sha256Encoder := sha256.New()
	sha256Encoder.Write([]byte(encPassword))
	encPassword = fmt.Sprintf("%x", sha256Encoder.Sum(nil))

	encPassword = fmt.Sprintf("%s%s%s", saltBeginning, encPassword, saltEnd)

	sha256Encoder.Write([]byte(encPassword))
	encPassword = fmt.Sprintf("%x", sha256Encoder.Sum(nil))
	return encPassword
}
