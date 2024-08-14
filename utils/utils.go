package utils

import (
	"strings"
)

func ExtractWildcardValues(pattern, filled string) []string {
	// Split both pattern and filled strings into segments
	patternSegments := strings.Split(pattern, ".")
	filledSegments := strings.Split(filled, ".")

	// Check if the number of segments matches
	if len(filledSegments) != len(patternSegments) {
		return nil
	}

	// Prepare a slice to hold the wildcard values
	wildcardValues := make([]string, 0, len(patternSegments))

	// Iterate through the segments
	for i := range patternSegments {
		switch patternSegments[i] {
		case "*":
			// If the segment is a wildcard, add the corresponding filled segment
			wildcardValues = append(wildcardValues, filledSegments[i])
		default:
			// If the segment is not a wildcard, check if it matches the filled segment
			if patternSegments[i] != filledSegments[i] {
				return nil
			}
		}
	}

	return wildcardValues
}
