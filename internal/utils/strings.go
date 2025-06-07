package utils

import "strings"

func GetFirstItem(input string, delimiter string) string {
	parts := strings.Split(input, delimiter)
	if len(parts) > 0 {
		return parts[0]
	}
	return input
}
