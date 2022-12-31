package helpers

import (
	"testing"

	"github.com/srackham/cryptor/internal/assert"
)

func TestParseDateOrOffset(t *testing.T) {
	type test struct {
		input string
		want  string
	}
	tests := []test{
		{input: "2022-12-31", want: "2022-12-31"},
		{input: "0", want: "2022-12-20"},
		{input: "-1", want: "2022-12-19"},
		{input: "-10", want: "2022-12-10"},
		{input: "1", want: "2022-12-21"},
	}
	for _, tc := range tests {
		got, err := ParseDateOrOffset(tc.input, "2022-12-20")
		assert.PassIf(t, err == nil, "error parsing date: %q", tc.input)
		assert.PassIf(t, tc.want == got, "input: %q: wanted: %q: got: %q", tc.input, tc.want, got)
	}
	_, err := ParseDateOrOffset("foo", "2022-12-20")
	assert.PassIf(t, err != nil, "expected error")
	_, err = ParseDateOrOffset("0", "bar")
	assert.PassIf(t, err != nil, "expected error")
}
