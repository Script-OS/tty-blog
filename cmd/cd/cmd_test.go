package cd

import (
	"testing"
	"tty-blog/global"
)

// TestCDPathCalculation tests the path calculation logic
func TestCDPathCalculation(t *testing.T) {
	tests := []struct {
		name     string
		workDir  string
		arg      string
		expected string
	}{
		{
			name:     "relative subdir",
			workDir:  "articles",
			arg:      "drafts",
			expected: "articles/drafts",
		},
		{
			name:     "navigate up",
			workDir:  "articles",
			arg:      "..",
			expected: ".",
		},
		{
			name:     "nested path",
			workDir:  ".",
			arg:      "articles/drafts",
			expected: "articles/drafts",
		},
		{
			name:     "absolute path",
			workDir:  "articles",
			arg:      "/root",
			expected: "root",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			global.WorkDir = tt.workDir
			result := global.CalcPath(tt.arg)
			if result != tt.expected {
				t.Errorf("CalcPath(%s) from %s = %s, expected %s",
					tt.arg, tt.workDir, result, tt.expected)
			}
		})
	}
}

// TestCDWorkDirUpdate tests that WorkDir is properly updated
func TestCDWorkDirUpdate(t *testing.T) {
	global.WorkDir = "articles"

	// Simulate cd to subdir
	global.WorkDir = global.CalcPath("drafts")
	if global.WorkDir != "articles/drafts" {
		t.Errorf("WorkDir should be 'articles/drafts', got '%s'", global.WorkDir)
	}

	// Simulate cd to parent (from articles/drafts back to articles)
	global.WorkDir = global.CalcPath("..")
	if global.WorkDir != "articles" {
		t.Errorf("WorkDir should be 'articles', got '%s'", global.WorkDir)
	}

	// Go back to root
	global.WorkDir = global.CalcPath("..")
	if global.WorkDir != "." {
		t.Errorf("WorkDir should be '.', got '%s'", global.WorkDir)
	}
}
