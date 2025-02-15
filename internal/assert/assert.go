/*
	Simple assertions package.
*/

package assert

import (
	"math"
	"regexp"
	"strings"
	"testing"
)

// PassIf fails and prints formatted message if not ok.
func PassIf(t *testing.T, ok bool, format string, args ...any) {
	t.Helper()
	if strings.Contains(format, "%s") {
		t.Fatalf("use '%%v' instead of '%%s' to report possibly nil values: '%s'", format)
	}
	if !ok {
		// I don't want errors to cascade, prefer to deal with them one by one.
		t.Fatalf(format, args...)
		// t.Errorf(format, args...)
	}
}

// FailIf fails and prints formatted message if not ok.
func FailIf(t *testing.T, notOk bool, format string, args ...any) {
	PassIf(t, !notOk, format, args...)
}

func Equal[T comparable](t *testing.T, wanted, got T) {
	t.Helper()
	PassIf(t, got == wanted, "\nwanted: %#v\ngot:    %#v", wanted, got)
}

func NotEqual[T comparable](t *testing.T, wanted, got T) {
	t.Helper()
	PassIf(t, wanted != got, "should not be %#v", got)
}

func True(t *testing.T, got bool) {
	t.Helper()
	PassIf(t, got, "should be true")
}

func False(t *testing.T, got bool) {
	t.Helper()
	PassIf(t, !got, "should be false")
}

func EqualValues[T comparable](t *testing.T, wanted, got []T) {
	t.Helper()
	PassIf(t, len(got) == len(wanted), "\nwanted: %#v\ngot:    %#v", wanted, got)
	for k := range got {
		PassIf(t, got[k] == wanted[k], "\nwanted: %#v\ngot:    %#v", wanted, got)
	}
}

func Panics(t *testing.T, f func()) {
	t.Helper()
	defer func() {
		recover()
	}()
	f()
	t.Error("should have panicked")
}

func Contains(t *testing.T, s, substr string) {
	t.Helper()
	PassIf(t, strings.Contains(s, substr), "%q does not contain %q", s, substr)
}

func NotContains(t *testing.T, s, substr string) {
	t.Helper()
	PassIf(t, !strings.Contains(s, substr), "%q contains %q", s, substr)
}

func ContainsPattern(t *testing.T, s, pattern string) {
	t.Helper()
	matched, _ := regexp.MatchString(pattern, s)
	PassIf(t, matched, "%q does not contain pattern %q", s, pattern)
}

func findMismatchContext(s1, s2 string, n int) (string, string) {
	minLen := int(math.Min(float64(len(s1)), float64(len(s2))))
	for i := 0; i < minLen; i++ {
		if s1[i] != s2[i] {
			// Determine start and end indices for substring extraction
			start := int(math.Max(float64(i-n), 0))
			end1 := int(math.Min(float64(i+n+1), float64(len(s1))))
			end2 := int(math.Min(float64(i+n+1), float64(len(s2))))

			return s1[start:end1], s2[start:end2]
		}
	}
	// If lengths differ, consider it a mismatch at the shorter string's end
	if len(s1) != len(s2) {
		start := int(math.Max(float64(minLen-n), 0))
		end1 := int(math.Min(float64(minLen+n), float64(len(s1))))
		end2 := int(math.Min(float64(minLen+n), float64(len(s2))))
		return s1[start:end1], s2[start:end2]
	}
	// If no mismatch found, return empty strings
	return "", ""
}

func EqualStrings(t *testing.T, wanted, got string) {
	t.Helper()
	if got == wanted {
		return
	}
	wanted, got = findMismatchContext(wanted, got, 10)
	PassIf(t, got == wanted, "\nwanted: ..%#v..\ngot:    ..%#v..", wanted, got)
}
