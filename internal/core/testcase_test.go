package core

import "testing"

func TestOutputMatch(t *testing.T) {
	tests := []struct {
		name     string
		actual   string
		expected string
		want     bool
	}{
		{
			name:     "both empty",
			actual:   "",
			expected: "",
			want:     true,
		},
		{
			name:     "exact match",
			actual:   "hello world",
			expected: "hello world",
			want:     true,
		},
		{
			name:     "trailing newline ignored",
			actual:   "42\n",
			expected: "42",
			want:     true,
		},
		{
			name:     "leading newline ignored",
			actual:   "\n42",
			expected: "42",
			want:     true,
		},
		{
			name:     "extra whitespace between tokens",
			actual:   "1  2  3",
			expected: "1 2 3",
			want:     true,
		},
		{
			name:     "multiline match",
			actual:   "1\n2\n3\n",
			expected: "1\n2\n3",
			want:     true,
		},
		{
			name:     "unequal spaced multiline match",
			actual:   "1\n2\n\n3",
			expected: "1\n2\n3",
			want:     true,
		},
		{
			name:     "multiline with extra spaces per line",
			actual:   "a  b\nc  d\n",
			expected: "a b\nc d",
			want:     true,
		},
		{
			name:     "different values",
			actual:   "1",
			expected: "2",
			want:     false,
		},
		{
			name:     "different line count",
			actual:   "1\n2",
			expected: "1",
			want:     false,
		},
		{
			name:     "actual has extra line",
			actual:   "1\n2\n3",
			expected: "1\n2",
			want:     false,
		},
		{
			name:     "expected has extra line",
			actual:   "1\n2",
			expected: "1\n2\n3",
			want:     false,
		},
		{
			name:     "token count mismatch on a line",
			actual:   "1 2 3",
			expected: "1 2",
			want:     false,
		},
		{
			name:     "actual empty expected non-empty",
			actual:   "",
			expected: "hello",
			want:     false,
		},
		{
			name:     "actual non-empty expected empty",
			actual:   "hello",
			expected: "",
			want:     false,
		},
		{
			name:     "unequal spaced multiline match",
			actual:   "1\n2\n\n3",
			expected: "1\n2\n3",
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := outputMatches(tt.actual, tt.expected)
			if got != tt.want {
				t.Errorf("outputMatches(%q, %q) = %v, want %v", tt.actual, tt.expected, got, tt.want)
			}
		})
	}
}
