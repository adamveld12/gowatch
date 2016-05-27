package gowatch

import (
	"os"
	"testing"
	"time"
)

func Test_defaultArtifactOutputName(t *testing.T) {
	cases := []struct {
		input1   string
		input2   string
		expected string
	}{
		{"github.com/adamveld12/test", "", "test"},
		{"github.com/adamveld12/test", "app", "app"},
	}

	for _, c := range cases {
		actual := defaultArtifactOutputName(c.input1, c.input2)
		if actual != c.expected {
			t.Errorf("expected %v - actual %v", c.expected, actual)
		}
	}
}

func Test_defaultWait(t *testing.T) {
	cases := []struct {
		input    time.Duration
		expected time.Duration
	}{
		{time.Millisecond * 740, time.Millisecond * 740},
		{time.Millisecond, time.Millisecond * 500},
		{time.Millisecond * 300, time.Millisecond * 500},
		{time.Second * 11, time.Second * 10},
		{time.Second * 9, time.Second * 9},
	}

	for _, c := range cases {
		actual := defaultWait(c.input)
		if actual != c.expected {
			t.Errorf("expected %v - actual %v", c.expected, actual)
		}
	}

}

func Test_defaultPackagePath(t *testing.T) {
	abs, _ := os.Getwd()
	cases := []struct {
		input    string
		expected string
	}{
		{".", abs},
		{"", abs},
		{"test/folder", abs + "/test/folder"},
	}

	for _, c := range cases {
		actual := defaultPackagePath(c.input)
		if actual != c.expected {
			t.Errorf("expected %v - actual %v", c.expected, actual)
		}
	}
}
