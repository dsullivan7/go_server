package server

import (
	"fmt"
	"crypto/cipher"
	"encoding/hex"
)

type Cipher struct {
	cipher cipher.Block
}

func NewCipher(cphr cipher.Block) Cipher {
  return Cipher{
    cipher: cphr,
  }
}

func (cphr *Cipher) Encrypt(phrase string) (string) {
  out := make([]byte, len(phrase))
  cphr.cipher.Encrypt(out, []byte(phrase))

	return hex.EncodeToString(out)
}

func (cphr *Cipher) Decrypt(phrase string) (string, error) {
    phraseDecoded, errDecode := hex.DecodeString(phrase)

    if (errDecode != nil) {
      return "", fmt.Errorf("error decoding hex: %w", errDecode)
    }

    phraseDecrypted := make([]byte, len(phraseDecoded))
    cphr.cipher.Decrypt(phraseDecrypted, phraseDecoded)

    return string(phraseDecrypted), nil
}
