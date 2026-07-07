package util

import "strings"

func TrimToLowerString(txt string) (string) {
	return strings.ToLower(strings.TrimSpace(txt))
}