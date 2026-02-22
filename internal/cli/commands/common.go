package commands

import (
	"regexp"
	"unicode/utf8"
)

var ansiEscapePattern = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func visibleRuneCount(s string) int {
	return utf8.RuneCountInString(ansiEscapePattern.ReplaceAllString(s, ""))
}
