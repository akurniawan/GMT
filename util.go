package gmt

import (
	"regexp"
	"unicode"
)

// Replacement is a tuple consisting of regex and their substitution
type Replacement struct {
	rgx *regexp.Regexp
	sub string
}

// NewReplacement creates Replacement object
func NewReplacement(rgx string, sub string) (replacement Replacement) {
	replacement = Replacement{regexp.MustCompile(rgx), sub}
	return
}

// Flatten will reduce dimensionality (2d to 1d) of the arguments
func Flatten(r [][]Replacement) []Replacement {
	var flattenedReplacement []Replacement

	for _, r0 := range r {
		for _, r1 := range r0 {
			flattenedReplacement = append(flattenedReplacement, r1)
		}
	}
	return flattenedReplacement
}

// IsLower checks whether all characters in a string are consisted of lowercase characters
func IsLower(text string) bool {
	for _, t := range text {
		res := unicode.IsLower(t)
		if !res {
			return false
		}
	}
	return true
}

// IsAnyAlphabet checks if alphabet character exist at least once in any string
func IsAnyAlphabet(text string) bool {
	totalMatch := 0
	for _, t := range text {
		res := unicode.IsLetter(t)
		if res {
			totalMatch++
		}
	}

	if totalMatch > 0 {
		return true
	}
	return false
}

// IsNumber checks whether all characters in a string are consisted of numbers
func IsNumber(text string) bool {
	for _, t := range text {
		if !unicode.IsNumber(t) {
			return false
		}
	}
	return true
}

// IsInArray checks if text is available in arr
func IsInArray(text string, arr []string) bool {
	for _, v := range arr {
		if text == v {
			return true
		}
	}
	return false
}

// RemoveEmptyStringFromSlice will check if any empty string exists on an array
// and remove them
func RemoveEmptyStringFromSlice(texts []string) (result []string) {
	for _, text := range texts {
		if text != "" {
			result = append(result, text)
		}
	}
	return
}
