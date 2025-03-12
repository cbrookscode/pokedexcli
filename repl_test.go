package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello    world   ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "  hello!! how? now browncow wow   ",
			expected: []string{"hello!!", "how?", "now", "browncow", "wow"},
		},
		{
			input:    "helloworld",
			expected: []string{"helloworld"},
		},
		{
			input:    "",
			expected: []string{},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("length of slice array doesn't equal length of expected slice. Expected: %d  Actual: %d, %v", len(c.expected), len(actual), actual)
			t.Fail()
		}
		for i := range actual {
			word := actual[i]
			expectedword := c.expected[i]
			if word != expectedword {
				t.Errorf("actual word and expected word don't match")
				t.Fail()
			}
		}
	}
}
