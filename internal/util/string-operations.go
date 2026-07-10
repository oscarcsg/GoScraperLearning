package util

import "strings"

func TrimToLowerString(txt string) (string) {
	return strings.ToLower(strings.TrimSpace(txt))
}

func CheckOnlyNumber(txt string) (bool) {
	if len(txt) == 0 {
		return false
	}

	for i := 0; i < len(txt); i++ {
		// Check every byte if it is between the char '0' (48) or '9' (57)
		if txt[i] < '0' || txt[i] > '9' {
			return false
		}
	}

	return true
}