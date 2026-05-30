package utilities

import (
	"strings"

	"github.com/agnivade/levenshtein"
)

var stopWords = map[string]bool{
	"the":  true,
	"and":  true,
	"or":   true,
	"a":    true,
	"an":   true,
	"of":   true,
	"to":   true,
	"in":   true,
	"on":   true,
	"for":  true,
	"with": true,
}

func GenerateSearchTokens(title string) []string {
	normalized := NormalizeText(title)

	words := strings.Fields(normalized)

	var searchTokens []string

	for _, word := range words {
		if !stopWords[word] {
			searchTokens = append(searchTokens, word)
		}
	}
	return searchTokens
}

func CalculateTokenScore(searchTokens, movieTokens []string) float64 {
	var score float64 = 0

	tokenMap := make(map[string]bool)

	for _, token := range movieTokens {
		tokenMap[token] = true
	}

	for _, token := range searchTokens {

		if tokenMap[token] {
			score++
		}
	}

	return score
}

func CalculateFuzzyTokenScore(searchTokens, movieTokens []string) float64 {
	totalScore := 0.0

	for _, searchToken := range searchTokens {

		bestMatch := 0.0

		for _, movieToken := range movieTokens {

			similarity := CalculateSimilarity(searchToken, movieToken)

			if similarity > bestMatch {
				bestMatch = similarity
			}
		}

		totalScore += bestMatch
	}

	return totalScore
}

func CalculateSimilarity(a, b string) float64 {
	if a == "" || b == "" {
		return 0
	}

	distance := levenshtein.ComputeDistance(a, b)

	maxLen := len(a)

	if len(b) > maxLen {
		maxLen = len(b)
	}

	return 1 - float64(distance)/float64(maxLen)
}
