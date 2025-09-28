package utils

import "fmt"

// GenerateSeatLabels generates seat labels like A1, A2, B1, B2...
func GenerateSeatLabels(total int) []string {
	labels := []string{}
	rows := []string{"A", "B", "C", "D", "E", "F", "G"}
	count := 0

	for _, row := range rows {
		for col := 1; col <= 10; col++ {
			labels = append(labels, fmt.Sprintf("%s%d", row, col))
			count++
			if count >= total {
				return labels
			}
		}
	}
	return labels
}

// Contains checks if a value exists in a slice
func Contains(slice []string, val string) bool {
	for _, s := range slice {
		if s == val {
			return true
		}
	}
	return false
}
