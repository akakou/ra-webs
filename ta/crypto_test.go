package ta

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSecureChannel(t *testing.T) {
	expected := []byte("hello")

	sc, err := newSecureChannel([]byte("01234567890123456789012345678901"))
	assert.NoError(t, err)

	cipher, err := sc.encrypt(expected)
	assert.NoError(t, err)

	actual, err := sc.decrypt(cipher)
	assert.NoError(t, err)

	assert.Equal(t, expected, actual)
}
