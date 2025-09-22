package shortener

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
)

const alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const alphabetLen = uint64(len(alphabet))

// generateRandomCode creates a pseudo-random base62 string of length n.
func GenerateCode(n int) (string, error) {
	if n <= 0 {
		return "", errors.New("length must be > 0")
	}
	b := make([]byte, n)
	// Use crypto/rand for unpredictability
	buf := make([]byte, 8)
	for i := 0; i < n; i++ {
		_, err := rand.Read(buf)
		if err != nil {
			return "", err
		}
		r := binary.LittleEndian.Uint64(buf)
		b[i] = alphabet[r%alphabetLen]
	}
	return string(b), nil
}
