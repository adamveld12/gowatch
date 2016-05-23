package main

import "testing"

// TestfmtIteration tests that strings are formtted properly
func Test_fmtIteration(t *testing.T) {

	cases := []struct {
		input    int
		expected string
	}{
		{0, "0 seconds and counting!\n"},
		{1, "1 second and counting!\n"},
	}

	for _, tc := range cases {
		actual := fmtIteration(tc.input)

		if actual != tc.expected {
			t.Errorf("expected %v - actual %v", tc.expected, actual)
		}
	}

}
