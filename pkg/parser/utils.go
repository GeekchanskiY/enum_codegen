package parser

import (
	"regexp"
	"strings"
)

// GetTranslationFromComment - get enum's translation from comment
func GetTranslationFromComment(comment string) string {
	re := regexp.MustCompile(`Translate="([^"]+)"`)
	match := re.FindStringSubmatch(comment)

	if len(match) > 1 {
		return match[1]
	} else {
		return ""
	}
}

// GetValueFromComment - get enum's string value from comment
func GetValueFromComment(comment string) string {
	re := regexp.MustCompile(`Value="([^"]+)"`)
	match := re.FindStringSubmatch(comment)

	if len(match) > 1 {
		return match[1]
	} else {
		return ""
	}
}

// CamelToSnake converts camelCase or PascalCase with numbers to snake_case
func CamelToSnake(s string) string {
	re1 := regexp.MustCompile("([a-z0-9])([A-Z])")
	s = re1.ReplaceAllString(s, "${1}_${2}")

	re2 := regexp.MustCompile("([a-zA-Z])([0-9])")
	s = re2.ReplaceAllString(s, "${1}_${2}")

	re3 := regexp.MustCompile("([0-9])([a-zA-Z])")
	s = re3.ReplaceAllString(s, "${1}_${2}")

	return strings.ToLower(s)
}
