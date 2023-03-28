package helpers

import (
	"fmt"
	"os/exec"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"golang.org/x/exp/constraints"
)

func If[T any](b bool, trueVal, falseVal T) T {
	if b {
		return trueVal
	} else {
		return falseVal
	}
}

// NormalizeNewlines converts \r\n (Windows) and \n (Mac OS) line terminations to
// \n (UNIX) termination.
func NormalizeNewlines(s string) (result string) {
	result = strings.ReplaceAll(s, "\r\n", "\n")
	result = strings.ReplaceAll(result, "\r", "\n")
	return
}

// SortedMapKeys returns a sorted array of map keys.
func SortedMapKeys[K constraints.Ordered, V any](m map[K]V) (res []K) {
	for k := range m {
		res = append(res, k)
	}
	sort.Slice(res, func(i, j int) bool { return res[i] < res[j] })
	return
}

// CopyMap returns a copy of a map.
func CopyMap[K comparable, V any](m map[K]V) (result map[K]V) {
	result = make(map[K]V)
	for k, v := range m {
		result[k] = v
	}
	return result
}

// MergeMap merges maps into dst map.
func MergeMap[K comparable, V any](dst map[K]V, maps ...map[K]V) {
	for _, m := range maps {
		for k, v := range m {
			dst[k] = v
		}
	}
}

// LaunchBrowser launches the browser at the url address. Waits till launch
// completed. Credit: https://stackoverflow.com/a/39324149/1136455
func LaunchBrowser(url string) error {
	var cmd string
	var args []string
	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Run()
}

// ParseDateString parses converts a date string to a time.Time. If timezone is not specified Local is assumed.
func ParseDateString(text string, loc *time.Location) (time.Time, error) {
	if loc == nil {
		loc, _ = time.LoadLocation("Local")
	}
	text = strings.TrimSpace(text)
	d, err := time.Parse(time.RFC3339, text)
	if err != nil {
		if d, err = time.Parse("2006-01-02 15:04:05-07:00", text); err != nil {
			if d, err = time.ParseInLocation("2006-01-02 15:04:05", text, loc); err != nil {
				d, err = time.ParseInLocation("2006-01-02", text, loc)
			}
		}
	}
	if err != nil {
		err = fmt.Errorf("invalid date value: \"%s\"", text)
	}
	return d, err
}

// TodaysDate returns the current date as a string.
func TodaysDate() string {
	return time.Now().Format("2006-01-02")
}

// IsDateString returns true if the `data` is formatted like YYYY-MM-DD.
func IsDateString(date string) bool {
	if _, err := time.Parse("2006-01-02", date); err != nil {
		return false
	} else {
		return true
	}
}

// TimeNowString returns the current time as a string.
func TimeNowString() string {
	return time.Now().Format("15:04:05")
}

// ParseDateOrOffset attempts to parse the `date` string which can be a YYYY-MM-DD
// formatted date or an integer offset (0=today, -1=yesterday, ...).
// The `fromDate` is a YYYY-MM-DD date string representing the date offset origin.
func ParseDateOrOffset(date string, fromDate string) (string, error) {
	if IsDateString(date) {
		return date, nil
	}
	var offset int
	var err error
	if offset, err = strconv.Atoi(date); err != nil {
		return "", fmt.Errorf("invalid date: \"%s\"", date)
	}
	d, err := time.Parse("2006-01-02", fromDate)
	if err != nil {
		return "", fmt.Errorf("invalid date: \"%s\"", fromDate)
	}
	d = d.AddDate(0, 0, offset)
	return d.Format("2006-01-02"), nil
}
