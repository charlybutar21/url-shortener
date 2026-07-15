package base62

import (
	"errors"
	"strings"
)

const alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const base = uint64(len(alphabet))

// Encode converts a uint64 ID to a Base62 string
func Encode(id uint64) string {
	if id == 0 {
		return string(alphabet[0])
	}

	var chars []byte
	for id > 0 {
		rem := id % base
		chars = append(chars, alphabet[rem])
		id = id / base
	}

	// Reverse the characters
	for i, j := 0, len(chars)-1; i < j; i, j = i+1, j-1 {
		chars[i], chars[j] = chars[j], chars[i]
	}

	return string(chars)
}

// Decode converts a Base62 string back to a uint64 ID
func Decode(encoded string) (uint64, error) {
	var id uint64

	for _, char := range encoded {
		pos := strings.IndexRune(alphabet, char)
		if pos == -1 {
			return 0, errors.New("invalid character in base62 string")
		}
		
		id = id*base + uint64(pos)
	}

	return id, nil
}
