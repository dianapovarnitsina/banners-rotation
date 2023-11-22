package logger

import (
	"bytes"
	"strings"
	"testing"
)

func TestLogger(t *testing.T) {
	tests := []struct {
		level    LogLevel
		expected string
	}{
		{Debug, "DEBUG "},
		{Info, "INFO "},
		{Warning, "WARNING "},
		{Error, "ERROR "},
	}

	for _, test := range tests {
		t.Run(string(test.level), func(t *testing.T) {
			var buf bytes.Buffer
			logger := New(string(test.level), &buf)

			// Log a message
			logger.Log("Test message")

			// Check if the logged message contains the expected level prefix
			if !strings.HasPrefix(buf.String(), test.expected) {
				t.Errorf("Expected message to start with '%s', but got: '%s'", test.expected, buf.String())
			}
		})
	}
}
