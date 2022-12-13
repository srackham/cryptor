package helpers

import (
	"fmt"
	"os/exec"
	"runtime"
	"sort"
	"strings"
	"time"

	"golang.org/x/exp/constraints"
)

// TODO helpers_test.go

func If[T any](b bool, trueVal, falseVal T) T {
	if b {
		return trueVal
	} else {
		return falseVal
	}
}

// NormalizeNewlines converts \r\n (Window) and \n (Mac OS) line terminations to
// \n (UNIX) termination.
func NormalizeNewlines(s string) (result string) {
	result = strings.ReplaceAll(s, "\r\n", "\n")
	result = strings.ReplaceAll(result, "\r", "\n")
	return
}

// SortedMapKeys returns a sorted array of map keys.
// TODO tests
func SortedMapKeys[K constraints.Ordered, V any](m map[K]V) (res []K) {
	for k := range m {
		res = append(res, k)
	}
	sort.Slice(res, func(i, j int) bool { return res[i] < res[j] })
	return
}

// CopyMap returns a copy of a map.
// TODO generalise keys
func CopyMap[T any](m map[string]T) (result map[string]T) {
	result = map[string]T{}
	for k, v := range m {
		result[k] = v
	}
	return
}

// MergeMap merges maps into dst map.
// TODO generalise keys
func MergeMap[T any](dst map[string]T, maps ...map[string]T) {
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
		err = fmt.Errorf("illegal date value: %q", text)
	}
	return d, err
}

// DateNowString returns the current date as a string.
func DateNowString() string {
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
