package renderer

import (
	"bytes"
	"testing"
)

// TestTableRender tests that table rendering does not panic with negative Repeat count
func TestTableRender(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name: "simple table",
			input: `| Name | Age |
|------|-----|
| Alice | 25 |
| Bob | 30 |`,
			wantErr: false,
		},
		{
			name: "table with long content",
			input: `| Short | VeryLongColumnName |
|-------|--------------------|
| A | B |
| ThisIsAVeryLongContentThatExceedsColumnWidth | Short |`,
			wantErr: false,
		},
		{
			name: "table with ANSI styled content simulation",
			input: `| Command | Description |
|---------|-------------|
| ls | List directory contents |
| cd | Change directory |
| view | View markdown file with rendering |`,
			wantErr: false,
		},
		{
			name: "empty table cells",
			input: `| A | B | C |
|---|---|---|
| 1 | | 3 |
| | 2 | |`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			md := New(80)
			var buf bytes.Buffer
			err := md.Convert([]byte(tt.input), &buf)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}

			// Ensure output is not empty for valid tables
			if !tt.wantErr && buf.Len() == 0 {
				t.Errorf("expected non-empty output")
			}
		})
	}
}

// TestRenderCellNegativePadding tests the fix for negative Repeat count
func TestRenderCellNegativePadding(t *testing.T) {
	ctx := &RenderContext{
		Meta: map[string]interface{}{
			"table": []int{5, 5, 5}, // Small column widths
			"row":   "",
			"col":   0,
		},
		Deco: []BlockDecorator{&BlankDeco{}},
	}

	// Create content that would exceed column width if bug not fixed
	content := "ThisIsLongContentExceedingWidth"
	actions := []Action{}

	// This should not panic due to negative Repeat count
	result := RenderCell(ctx, nil, 80, content, actions, nil)

	if result != "" {
		t.Errorf("RenderCell should return empty string, got: %s", result)
	}

	// Check that row was updated (with content, no negative padding crash)
	row := ctx.Meta["row"].(string)
	if row == "" {
		t.Errorf("row should be updated after RenderCell")
	}
}

// TestRenderCellNormalPadding tests normal padding behavior
func TestRenderCellNormalPadding(t *testing.T) {
	ctx := &RenderContext{
		Meta: map[string]interface{}{
			"table": []int{20, 20, 20}, // Large column widths
			"row":   "",
			"col":   0,
		},
		Deco: []BlockDecorator{&BlankDeco{}},
	}

	// Content shorter than column width - should be padded
	content := "Short"
	actions := []Action{}

	result := RenderCell(ctx, nil, 80, content, actions, nil)

	if result != "" {
		t.Errorf("RenderCell should return empty string, got: %s", result)
	}

	// Content should be padded to column width
	row := ctx.Meta["row"].(string)
	if row == "" {
		t.Errorf("row should be updated after RenderCell")
	}
}

// TestEasyRenderTable tests the full rendering pipeline with tables
func TestEasyRenderTable(t *testing.T) {
	md := New(80)

	markdown := `# Test Document

| Column A | Column B |
|----------|----------|
| Value 1 | Value 2 |
| Value 3 | Value 4 |

Some text after table.
`

	result, err := EasyRender(md, []byte(markdown))
	if err != nil {
		t.Errorf("EasyRender failed: %v", err)
	}

	if len(result) == 0 {
		t.Errorf("expected non-empty result")
	}
}
