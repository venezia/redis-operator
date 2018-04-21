package main

import (
	"testing"
	"os"
)

func TestHomeDir(t *testing.T) {
	expected := "failure"
	if os.Getenv("HOME") != "" {
		expected = os.Getenv("HOME")
	} else {
		expected = os.Getenv("USERPROFILE")
	}
	actual := homeDir()
	if actual != expected {
		t.Errorf("Test failed, got -->%s<-- was expecting -->%s<--", actual, expected)
	}
}

