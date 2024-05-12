package text

import (
	"errors"
	"unicode"
)

func isAlphanumeric(s string) bool {
	for _, char := range s {
		if !unicode.IsLower(char) && !unicode.IsDigit(char) {
			return false
		}
	}
	return true
}

var (
	ErrInvalidLength   = errors.New("text: invalid length")
	ErrFirstLetter     = errors.New("text: first letter is not a letter")
	ErrNotAlphaNumeric = errors.New("text: not alphanumeric")
)
