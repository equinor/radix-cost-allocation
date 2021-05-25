package slice

import "strings"

func ContainsString(slice []string, value string) bool {
	for _, s := range slice {
		if s == value {
			return true
		}
	}

	return false
}

func ToLowerCase(slice []string) []string {
	lowerSlice := make([]string, len(slice), len(slice))

	for i, s := range slice {
		lowerSlice[i] = strings.ToLower(s)
	}

	return lowerSlice
}
