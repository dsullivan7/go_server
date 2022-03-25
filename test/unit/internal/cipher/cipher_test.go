package alpaca_test

import (
	"go_server/internal/cipher"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncrypt(t *testing.T) {
	t.Parallel()

	// needs to be 32 bytes long
	key := "01234567890123456789012345678901"

	testString := "testString"

	cphr := cipher.NewCipher()

	encrypted, errEncrypt := cphr.Encrypt(testString, key)

	assert.Nil(t, errEncrypt)

	assert.NotNil(t, encrypted)

	decrypted, errDecrypt := cphr.Decrypt(encrypted, key)
	assert.Nil(t, errDecrypt)
	assert.Equal(t, decrypted, "testString")
}

func TestDecrypt(t *testing.T) {
	t.Parallel()

	// needs to be 32 bytes long
	key := "01234567890123456789012345678901"

	cphr := cipher.NewCipher()

	testString := "a474f78e24e415e233d7e3db656cbd54a04f787efa1ce696cdf29e5ca2299e797334362c862b"

	decrypted, errDecrypt := cphr.Decrypt(testString, key)

	assert.Nil(t, errDecrypt)

	assert.Equal(t, decrypted, "testString")
}
