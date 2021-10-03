package stringfn

import "unicode"

// IsOnlyLetters check if a string content only letters
func IsOnlyLetters(text string) bool {
	for _, l:= range text {
		if !unicode.IsLetter(l) {
			return false
		}
	}

	return true
}
