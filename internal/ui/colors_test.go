package ui

import (
	"strings"
	"testing"
)

func TestSuccess(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantSymbol string
		wantColor  string
	}{
		{
			name:       "success with simple message",
			input:      "deployment complete",
			wantSymbol: checkMark,
			wantColor:  colorGreen,
		},
		{
			name:       "success with empty message",
			input:      "",
			wantSymbol: checkMark,
			wantColor:  colorGreen,
		},
		{
			name:       "success with special chars",
			input:      "/srv/myapp deployed",
			wantSymbol: checkMark,
			wantColor:  colorGreen,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Success(tt.input)

			if !strings.Contains(result, tt.wantColor) {
				t.Errorf("Success() missing color code")
			}
			if !strings.Contains(result, tt.wantSymbol) {
				t.Errorf("Success() missing checkmark")
			}
			if !strings.Contains(result, tt.input) {
				t.Errorf("Success() missing input message")
			}
			if !strings.Contains(result, colorReset) {
				t.Errorf("Success() missing color reset")
			}
		})
	}
}

func TestError(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"error with simple message", "deployment failed"},
		{"error with empty message", ""},
		{"error with path", "/etc/myapp permission denied"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Error(tt.input)

			if !strings.Contains(result, colorRed) {
				t.Errorf("Error() missing red color")
			}
			if !strings.Contains(result, cross) {
				t.Errorf("Error() missing cross symbol")
			}
			if !strings.Contains(result, tt.input) {
				t.Errorf("Error() missing input message")
			}
			if !strings.Contains(result, colorReset) {
				t.Errorf("Error() missing color reset")
			}
		})
	}
}

func TestWorking(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"working with simple message", "validating layout"},
		{"working with empty message", ""},
		{"working with status", "extracting assets"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Working(tt.input)

			if !strings.Contains(result, colorBlue) {
				t.Errorf("Working() missing blue color")
			}
			if !strings.Contains(result, arrow) {
				t.Errorf("Working() missing arrow symbol")
			}
			if !strings.Contains(result, tt.input) {
				t.Errorf("Working() missing input message")
			}
			if !strings.Contains(result, colorReset) {
				t.Errorf("Working() missing color reset")
			}
		})
	}
}

func TestWarn(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"warn with simple message", "permission issue detected"},
		{"warn with empty message", ""},
		{"warn with context", "file permissions differ from source"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Warn(tt.input)

			if !strings.Contains(result, colorYellow) {
				t.Errorf("Warn() missing yellow color")
			}
			if !strings.Contains(result, "⚠") {
				t.Errorf("Warn() missing warning symbol")
			}
			if !strings.Contains(result, tt.input) {
				t.Errorf("Warn() missing input message")
			}
			if !strings.Contains(result, colorReset) {
				t.Errorf("Warn() missing color reset")
			}
		})
	}
}

func TestColorConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant string
		want     string
	}{
		{"ColorReset", colorReset, "\033[0m"},
		{"ColorGreen", colorGreen, "\033[32m"},
		{"ColorRed", colorRed, "\033[31m"},
		{"ColorYellow", colorYellow, "\033[33m"},
		{"ColorBlue", colorBlue, "\033[34m"},
		{"ColorGray", colorGray, "\033[90m"},
		{"CheckMark", checkMark, "✓"},
		{"Cross", cross, "✗"},
		{"Arrow", arrow, "→"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.want {
				t.Errorf("constant mismatch: got %q, want %q", tt.constant, tt.want)
			}
		})
	}
}

func TestOutputFormat(t *testing.T) {
	// Verify exact format: color + symbol + reset + space + message
	result := Success("test")
	expected := colorGreen + checkMark + " " + "test" + colorReset

	if result != expected {
		t.Errorf("Success() format mismatch\ngot:  %q\nwant: %q", result, expected)
	}
}
