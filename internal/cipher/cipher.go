package cipher

import (
	"fmt"
	"io"
	"crypto/rand"
	"crypto/cipher"
	"crypto/aes"
	"encoding/hex"
)

type ICipher interface {
  Encrypt(phrase string, key string) (string, error)
  Decrypt(phrase string, key string) (string, error)
}

type Cipher struct {}

func NewCipher() Cipher {
  return Cipher{}
}

func (cphr *Cipher) Encrypt(phrase string, key string) (string, error) {
	keyBytes := []byte(key)
	aesCipher, errCipher := aes.NewCipher(keyBytes)

	if errCipher != nil {
		return "", fmt.Errorf("error creating cipher: %w", errCipher)
	}

	gcm, errGCM := cipher.NewGCM(aesCipher)

	if errGCM != nil {
		return "", fmt.Errorf("error creating gcm: %w", errGCM)
	}

	nonce := make([]byte, gcm.NonceSize())

	// populate nonce with a cryptographically secure random sequence
	_, errNonce := io.ReadFull(rand.Reader, nonce)

	if errNonce != nil {
		return "", fmt.Errorf("error populating nonce: %w", errNonce)
	}

	encryptedText := gcm.Seal(nonce, nonce, []byte(phrase), nil)

	return fmt.Sprintf("%x", encryptedText), nil
}

func (cphr *Cipher) Decrypt(phrase string, key string) (string, error) {
	keyBytes := []byte(key)
	aesCipher, errCipher := aes.NewCipher(keyBytes)

	if errCipher != nil {
		return "", fmt.Errorf("error creating cipher: %w", errCipher)
	}

	phraseBytes, errPhraseDecode := hex.DecodeString(phrase)

	if errPhraseDecode != nil {
		return "", fmt.Errorf("error decoding phrase string: %w", errPhraseDecode)
	}

	aesGCM, errGCM := cipher.NewGCM(aesCipher)
	if errGCM != nil {
		return "", fmt.Errorf("error creating gcm: %w", errGCM)
	}

	nonceSize := aesGCM.NonceSize()

	nonce, phraseCipher := phraseBytes[:nonceSize], phraseBytes[nonceSize:]

	phraseDecrypted, errDecrypt := aesGCM.Open(nil, nonce, phraseCipher, nil)

	if errDecrypt != nil {
		return "", fmt.Errorf("error decrypting: %w", errDecrypt)
	}

	return string(phraseDecrypted), nil
}
