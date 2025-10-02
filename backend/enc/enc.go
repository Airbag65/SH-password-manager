package enc

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"log"
)


func GenerateKeys() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}
	publicKey := privateKey.PublicKey
	return privateKey, &publicKey, nil
}

func PrivateKeyToPemString(key *rsa.PrivateKey) string {
	keyDER := x509.MarshalPKCS1PrivateKey(key)
	keyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: keyDER,
	})
	return string(keyPEM)
}

func PemStringToPrivateKey(pemString string) *rsa.PrivateKey {
	block, _ := pem.Decode([]byte(pemString))
	parseRes, _ := x509.ParsePKCS1PrivateKey(block.Bytes)
	return parseRes
}

func PublicKeyToPemString(key *rsa.PublicKey) string {
	keyDER := x509.MarshalPKCS1PublicKey(key)
	keyPEM := pem.EncodeToMemory(&pem.Block{
		Type: "PUBLIC KEY",
		Bytes: keyDER,
	})
	return string(keyPEM)
}

func PemStringToPublicKey(pemString string) *rsa.PublicKey {
	block, _ := pem.Decode([]byte(pemString))
	parseRes, _ := x509.ParsePKCS1PublicKey(block.Bytes)
	return parseRes
}


func Encrypt(message string, publicKey *rsa.PublicKey) ([]byte, error) {
	encBytes, err := rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		publicKey,
		[]byte(message),
		nil)
	if err != nil {
		log.Println("Could not Encrypt")
		return []byte(""), err
	}
	return encBytes, nil
}

func Decrypt(encryptedMessage []byte, privateKey *rsa.PrivateKey) (string, error) {
	decMessage, err := privateKey.Decrypt(nil, encryptedMessage, &rsa.OAEPOptions{Hash: crypto.SHA256})
	if err != nil {
		log.Println("Could not Decrypt")
		return "", err
	}
	return string(decMessage), nil
}
