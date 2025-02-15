package helpers

import (
	"fmt"
	"os"
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
	// Test with a map[string]string
	t.Run("CopiesStringStringMap", func(t *testing.T) {
		m := map[string]string{
			"a": "A",
			"b": "B",
			"c": "C",
		}

		result := CopyMap(m)

		if len(result) != len(m) {
			t.Errorf("Expected length of %d, but got %d", len(m), len(result))
		}

		for k, v := range result {
			if m[k] != v {
				t.Errorf("Expected value %s for key %s, but was %s", m[k], k, v)
			}
		}
	})

	// Test with a map[int]string
	t.Run("CopiesIntStringMap", func(t *testing.T) {
		m := map[int]string{
			1: "One",
			2: "Two",
			3: "Three",
		}

		result := CopyMap(m)

		if len(result) != len(m) {
			t.Errorf("Expected length of %d, but got %d", len(m), len(result))
		}

		for k, v := range result {
			if m[k] != v {
				t.Errorf("Expected value %s for key %v, but was %s", m[k], k, v)
			}
		}
	})
}

func TestMergeMap(t *testing.T) {
	// Test with a map[string]string
	t.Run("MergesStringStringMaps", func(t *testing.T) {
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

		for k, v := range expected {
			if dst[k] != v {
				t.Errorf("Expected value %s for key %s, but was %s", v, k, dst[k])
			}
		}
	})

	// Test with a map[int]string
	t.Run("MergesIntStringMaps", func(t *testing.T) {
		dst := map[int]string{
			1: "One",
			2: "Two",
			3: "Three",
		}

		m1 := map[int]string{
			4: "Four",
			5: "Five",
			6: "Six",
		}

		m2 := map[int]string{
			7: "Seven",
			8: "Eight",
			9: "Nine",
		}

		MergeMap(dst, m1, m2)

		expected := map[int]string{
			1: "One",
			2: "Two",
			3: "Three",
			4: "Four",
			5: "Five",
			6: "Six",
			7: "Seven",
			8: "Eight",
			9: "Nine",
		}

		for k, v := range expected {
			if dst[k] != v {
				t.Errorf("Expected value %s for key %d, but was %s", v, k, dst[k])
			}
		}
	})
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

func TestGetXDGConfigDir(t *testing.T) {
	// Save the original environment variables
	originalXDGConfigHome := os.Getenv("XDG_CONFIG_HOME")
	originalHome := os.Getenv("HOME")
	defer func() {
		os.Setenv("XDG_CONFIG_HOME", originalXDGConfigHome)
		os.Setenv("HOME", originalHome)
	}()

	testCases := []struct {
		name           string
		xdgConfigHome  string
		home           string
		expectedResult string
	}{
		{
			name:           "XDG_CONFIG_HOME set",
			xdgConfigHome:  "/custom/config/path",
			home:           "/home/user",
			expectedResult: "/custom/config/path",
		},
		{
			name:           "XDG_CONFIG_HOME not set",
			xdgConfigHome:  "",
			home:           "/home/user",
			expectedResult: "/home/user/.config",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			os.Setenv("XDG_CONFIG_HOME", tc.xdgConfigHome)
			os.Setenv("HOME", tc.home)
			result := GetXDGConfigDir()
			if result != tc.expectedResult {
				t.Errorf("Expected %s, but got %s", tc.expectedResult, result)
			}
		})
	}
}

func TestGetXDGCacheDir(t *testing.T) {
	// Save the original environment variables
	originalXDGCacheHome := os.Getenv("XDG_CACHE_HOME")
	originalHome := os.Getenv("HOME")
	defer func() {
		os.Setenv("XDG_CACHE_HOME", originalXDGCacheHome)
		os.Setenv("HOME", originalHome)
	}()

	testCases := []struct {
		name           string
		xdgCacheHome   string
		home           string
		expectedResult string
	}{
		{
			name:           "XDG_CACHE_HOME set",
			xdgCacheHome:   "/custom/cache/path",
			home:           "/home/user",
			expectedResult: "/custom/cache/path",
		},
		{
			name:           "XDG_CACHE_HOME not set",
			xdgCacheHome:   "",
			home:           "/home/user",
			expectedResult: "/home/user/.cache",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			os.Setenv("XDG_CACHE_HOME", tc.xdgCacheHome)
			os.Setenv("HOME", tc.home)
			result := GetXDGCacheDir()
			if result != tc.expectedResult {
				t.Errorf("Expected %s, but got %s", tc.expectedResult, result)
			}
		})
	}
}
