package strutil

import (
	"strings"
	"unicode/utf8"
)

// StringIn check if item is in ss
func StringIn(ss []string, item string) bool {
	for _, s := range ss {
		if s == item {
			return true
		}
	}
	return false
}

// Cutn cut s to n bytes
func Cutn(s string, n int) string {
	if len(s) >= n {
		return s[:n]
	} else {
		return s
	}
}

// CutnRune cut s to n runes
func CutnRune(s string, n int) string {
	if utf8.RuneCountInString(s) >= n {
		return string([]rune(s)[:n])
	} else {
		return s
	}
}

// ContainAny check if s contain any item in ss
func ContainAny(s string, ss []string) bool {
	for _, item := range ss {
		if strings.Contains(s, item) {
			return true
		}
	}

	return false
}

func EqualAny(s string, ss []string) bool {
	for _, item := range ss {
		if s == item {
			return true
		}
	}
	return false
}

func EncodeRuneToStr(r rune) string {
	p := make([]byte, 4)
	n := utf8.EncodeRune(p, r)
	return string(p[:n])
}

// reverser string
func ReverseString(s string) string {
	runes := []rune(s)
	for from, to := 0, len(runes)-1; from < to; from, to = from+1, to-1 {
		runes[from], runes[to] = runes[to], runes[from]
	}
	return string(runes)
}
