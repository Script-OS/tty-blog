package global

import (
	"testing"
)

// TestCalcPathRelative tests relative path calculation
func TestCalcPathRelative(t *testing.T) {
	// Set initial work directory
	WorkDir = "/articles"

	tests := []struct {
		input    string
		expected string
	}{
		{"drafts", "articles/drafts"},
		{"../posts", "posts"},
		{"./current", "articles/current"},
		{"subdir/file.md", "articles/subdir/file.md"},
	}

	for _, tt := range tests {
		result := CalcPath(tt.input)
		if result != tt.expected {
			t.Errorf("CalcPath(%s) = %s, expected %s", tt.input, result, tt.expected)
		}
	}
}

// TestCalcPathAbsolute tests absolute path handling
func TestCalcPathAbsolute(t *testing.T) {
	WorkDir = "/articles"

	tests := []struct {
		input    string
		expected string
	}{
		{"/root", "root"},
		{"/articles/drafts", "articles/drafts"},
		{"/", "."},
	}

	for _, tt := range tests {
		result := CalcPath(tt.input)
		if result != tt.expected {
			t.Errorf("CalcPath(%s) = %s, expected %s", tt.input, result, tt.expected)
		}
	}
}

// TestCalcPathClean tests path cleaning
func TestCalcPathClean(t *testing.T) {
	WorkDir = "/articles"

	tests := []struct {
		input    string
		expected string
	}{
		{"./drafts/../posts", "articles/posts"},
		{"drafts/./file.md", "articles/drafts/file.md"},
		{"../", "."},
	}

	for _, tt := range tests {
		result := CalcPath(tt.input)
		if result != tt.expected {
			t.Errorf("CalcPath(%s) = %s, expected %s", tt.input, result, tt.expected)
		}
	}
}

// TestCalcPathFromRoot tests path from root directory
func TestCalcPathFromRoot(t *testing.T) {
	WorkDir = "."

	tests := []struct {
		input    string
		expected string
	}{
		{"articles", "articles"},
		{"./file.md", "file.md"},
		{"subdir", "subdir"},
	}

	for _, tt := range tests {
		result := CalcPath(tt.input)
		if result != tt.expected {
			t.Errorf("CalcPath(%s) = %s, expected %s", tt.input, result, tt.expected)
		}
	}
}
