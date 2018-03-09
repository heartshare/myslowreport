package utils

import (
	"strings"
)

func Substr(str string, start, length int) string {
    rs := []rune(str)
    rl := len(rs)
    end := 0

    if start < 0 {
        start = rl - 1 + start
    }
    end = start + length

    if start > end {
        start, end = end, start
    }

    if start < 0 {
        start = 0
    }
    if start > rl {
        start = rl
    }
    if end < 0 {
        end = 0
    }
    if end > rl {
        end = rl
    }

    return string(rs[start:end])
}

func UnicodeIndex(str, substr string) int {
	result := strings.Index(str,substr)
	if result >= 0 {
		prefix := []byte(str)[0:result]
		rs := []rune(string(prefix))
		result = len(rs)
	}
	return result
}

func GetValue(str string) string {
	idx := UnicodeIndex(str, ":")
	val := Substr(str, idx + 1, len(str) - idx)
	val = strings.TrimLeft(val, " ")
	val = strings.TrimRight(val, " ")

	return val
}
