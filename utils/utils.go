package utils

import (
	"fmt"
	"strings"
)

func Contains[T comparable](value T, arr []T) bool {
	for _, k := range arr {
		if k == value {
			return true
		}
	}
	return false
}

// GetItemsMessage returns a message with the following format:
// message:
// + item1
// + item2
// + itemN
func GetItemsMessage(message string, items []string) string {
	formattedMessage := fmt.Sprintf("%s:\n\t+ ", message)
	formattedItems := strings.Join(items, "\n\t+ ")
	formattedItems += "\n"
	return formattedMessage + formattedItems
}
