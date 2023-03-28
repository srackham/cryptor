package helpers

import (
	"fmt"
	"testing"
	"time"
)

func TestIf(t *testing.T) {
	trueVal := "true"
	falseVal := "false"
	if res := If(true, trueVal, falseVal); res != trueVal {
		t.Errorf("Expected %s, got %s", trueVal, res)
	}
	if res := If(false, trueVal, falseVal); res != falseVal {
		t.Errorf("Expected %s, got %s", falseVal, res)
	}
}

func TestNormalizeNewlines(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"\r\n", "\n"},
		{"\r", "\n"},
		{"\n", "\n"},
		{"\r\n\r\n", "\n\n"},
		{"\r\n\r", "\n\n"},
		{"\r\r\n", "\n\n"},
	}
	for _, tc := range testCases {
		res := NormalizeNewlines(tc.input)
		if res != tc.expected {
			t.Errorf("Expected %s, got %s", tc.expected, res)
		}
	}
}

func TestSortedMapKeys(t *testing.T) {
	m := map[int]string{
		3: "three",
		1: "one",
		2: "two",
	}
	expected := []int{1, 2, 3}
	res := SortedMapKeys(m)
	if len(res) != len(expected) {
		t.Errorf("Expected %v, got %v", expected, res)
	}
	for i, v := range res {
		if v != expected[i] {
			t.Errorf("Expected %v, got %v", expected, res)
		}
	}
}

func TestCopyMap(t *testing.T) {
	m := map[string]string{
		"a": "A",
		"b": "B",
		"c": "C",
	}
	res := CopyMap(m)
	if len(res) != len(m) {
		t.Errorf("Expected %v, got %v", m, res)
	}
	for k, v := range res {
		if m[k] != v {
			t.Errorf("Expected %v, got %v", m, res)
		}
	}
}

func TestMergeMap(t *testing.T) {
	dst := map[string]string{
		"a": "A",
		"b": "B",
		"c": "C",
	}
	m1 := map[string]string{
		"d": "D",
		"e": "E",
		"f": "F",
	}
	m2 := map[string]string{
		"g": "G",
		"h": "H",
		"i": "I",
	}
	MergeMap(dst, m1, m2)
	expected := map[string]string{
		"a": "A",
		"b": "B",
		"c": "C",
		"d": "D",
		"e": "E",
		"f": "F",
		"g": "G",
		"h": "H",
		"i": "I",
	}
	if len(dst) != len(expected) {
		t.Errorf("Expected %v, got %v", expected, dst)
	}
	for k, v := range dst {
		if expected[k] != v {
			t.Errorf("Expected %v, got %v", expected, dst)
		}
	}
}

/* This test works only if you test with in an environment with a browser installed.
func TestLaunchBrowser(t *testing.T) {
	url := "https://www.example.com"
	err := LaunchBrowser(url)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}
*/

func TestParseDateString(t *testing.T) {
	loc, _ := time.LoadLocation("Local")
	testCases := []struct {
		input    string
		expected time.Time
	}{
		{"2020-01-02T15:04:05Z", time.Date(2020, 1, 2, 15, 4, 5, 0, time.UTC)},
		{"2020-01-02 15:04:05-07:00", time.Date(2020, 1, 2, 15, 4, 5, 0, time.FixedZone("", -7*60*60))},
		{"2020-01-02 15:04:05", time.Date(2020, 1, 2, 15, 4, 5, 0, loc)},
		{"2020-01-02", time.Date(2020, 1, 2, 0, 0, 0, 0, loc)},
	}
	for _, tc := range testCases {
		res, err := ParseDateString(tc.input, loc)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if !res.Equal(tc.expected) {
			t.Errorf("Expected %v, got %v", tc.expected, res)
		}
	}
}

func TestTodaysDate(t *testing.T) {
	now := time.Now()
	res := TodaysDate()
	expected := now.Format("2006-01-02")
	if res != expected {
		t.Errorf("Expected %s, got %s", expected, res)
	}
}

func TestIsDateString(t *testing.T) {
	testCases := []struct {
		input    string
		expected bool
	}{
		{"2020-01-02", true},
		{"2020-13-02", false},
		{"2020-01-32", false},
		{"2020-01-02 15:04:05", false},
	}
	for _, tc := range testCases {
		res := IsDateString(tc.input)
		if res != tc.expected {
			t.Errorf("Expected %v, got %v", tc.expected, res)
		}
	}
}

/* This naive test won't always work with this impure (stateful) function.
func TestTimeNowString(t *testing.T) {
	now := time.Now()
	res := TimeNowString()
	expected := now.Format("15:04:05")
	if res != expected {
		t.Errorf("Expected %s, got %s", expected, res)
	}
}
*/

func TestParseDateOrOffset(t *testing.T) {
	testCases := []struct {
		name          string
		date          string
		fromDate      string
		expected      string
		expectedError error
	}{
		{
			name:     "valid YYYY-MM-DD date",
			date:     "2021-06-01",
			fromDate: "2021-06-01",
			expected: "2021-06-01",
		},
		{
			name:     "valid integer offset (-1)",
			date:     "-1",
			fromDate: "2021-06-01",
			expected: "2021-05-31",
		},
		{
			name:     "valid integer offset (-10)",
			date:     "-10",
			fromDate: "2022-12-20",
			expected: "2022-12-10",
		},
		{
			name:          "invalid date string (month=13)",
			date:          "2021-13-01",
			fromDate:      "2021-06-01",
			expectedError: fmt.Errorf("invalid date: \"2021-13-01\""),
		},
		{
			name:          "invalid offset string (not an int)",
			date:          "not an int",
			fromDate:      "2021-06-01",
			expectedError: fmt.Errorf("invalid date: \"not an int\""),
		},
		{
			name:          "invalid fromDate string (bad format)",
			date:          "0",
			fromDate:      "20210601", // missing hyphens
			expectedError: fmt.Errorf("invalid date: \"20210601\""),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := ParseDateOrOffset(tc.date, tc.fromDate)

			if tc.expectedError != nil {
				// Ensure an error was returned
				if err == nil {
					t.Error("Expected error but got nil")
					return
				}

				// Ensure the error matches our expected error
				if err.Error() != tc.expectedError.Error() {
					t.Errorf("Expected error %q but got %q", tc.expectedError.Error(), err.Error())
				}
			} else {
				// Ensure no error was returned
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
					return
				}

				// Ensure the actual date matches our expected date
				if actual != tc.expected {
					t.Errorf("Expected %q but got %q", tc.expected, actual)
				}
			}
		})
	}
}
