package slice

import "strings"

func ToLowerCase(slice []string) []string {
	lowerSlice := make([]string, len(slice))

	for i, s := range slice {
		lowerSlice[i] = strings.ToLower(s)
	}

	return lowerSlice
}
