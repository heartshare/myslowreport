package utils

import (
	"strings"
	"os"
	"time"
	"strconv"
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

func FileExists(file string) bool {
	_, err := os.Stat(file)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func WeekdayCNString(t time.Time) string {
	weekDay := t.Weekday().String()
	if strings.Contains(weekDay, "Sunday") {
		return "星期天"
	} else if strings.Contains(weekDay, "Monday") {
		return "星期一"
	} else if strings.Contains(weekDay, "Tuesday") {
		return "星期二"
	} else if strings.Contains(weekDay, "Wednesday") {
		return "星期三"
	} else if strings.Contains(weekDay, "Thursday") {
		return "星期四"
	} else if strings.Contains(weekDay, "Friday") {
		return "星期五"
	} else if strings.Contains(weekDay, "Saturday") {
		return "星期六"
	} else {
		return ""
	}
}

//2006-01-02 15:04:05
func Today() time.Time {
	return time.Now()
}
func TodayStringByFormat(format string) string {
	return Today().Format(format)
}

func Yesterday() time.Time {
	return time.Now().AddDate(0, 0, -1)
}

func YesterdayString() string {
	return Yesterday().Format("20060102")
}

func YesterdayStringByFormat(format string) string {
	return Yesterday().Format(format)
}

func BeforeYesterday() time.Time {
	return time.Now().AddDate(0, 0, -2)
}

func BeforeYesterdayStringByFormat(format string) string {
	return BeforeYesterday().Format(format)
}

func StringToTimeByFormat(s string, format string) time.Time {
	t, _ := time.Parse(format, s)
	return t
}

func StringToFloat64(s string, defaultValue float64) float64 {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return defaultValue
	}
	return v
}