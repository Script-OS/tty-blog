package ls

import (
	"io/fs"
	"sort"
	"testing"
	"time"
)

// TestEntrySorter tests the sorting logic for directory entries
func TestEntrySorter(t *testing.T) {
	// Create mock directory entries
	now := time.Now()

	entries := EntrySorter{
		// File entries with different modification times
		createMockDirEntry("file1.md", false, now.Add(-1*time.Hour)),
		createMockDirEntry("file2.md", false, now.Add(-2*time.Hour)),
		createMockDirEntry("file3.md", false, now),
		// Directory entries
		createMockDirEntry("dir1", true, now.Add(-3*time.Hour)),
		createMockDirEntry("dir2", true, now.Add(-4*time.Hour)),
	}

	sort.Sort(entries)

	// Verify directories come first
	if !entries[0].IsDir() {
		t.Error("First entry should be a directory")
	}
	if !entries[1].IsDir() {
		t.Error("Second entry should be a directory")
	}

	// Verify files come after directories
	for i := 2; i < len(entries); i++ {
		if entries[i].IsDir() {
			t.Errorf("Entry %d should be a file, not directory", i)
		}
	}
}

// TestEntrySorterByTime tests that entries are sorted by modification time
func TestEntrySorterByTime(t *testing.T) {
	now := time.Now()

	entries := EntrySorter{
		createMockDirEntry("old.md", false, now.Add(-10*time.Hour)),
		createMockDirEntry("new.md", false, now),
		createMockDirEntry("middle.md", false, now.Add(-5*time.Hour)),
	}

	sort.Sort(entries)

	// Files should be sorted newest first
	if entries[0].Name() != "new.md" {
		t.Errorf("First file should be 'new.md', got '%s'", entries[0].Name())
	}
	if entries[1].Name() != "middle.md" {
		t.Errorf("Second file should be 'middle.md', got '%s'", entries[1].Name())
	}
	if entries[2].Name() != "old.md" {
		t.Errorf("Third file should be 'old.md', got '%s'", entries[2].Name())
	}
}

// TestEntrySorterEmpty tests sorting empty entries
func TestEntrySorterEmpty(t *testing.T) {
	entries := EntrySorter{}

	sort.Sort(entries)

	if len(entries) != 0 {
		t.Error("Empty sorter should remain empty")
	}
}

// TestEntrySorterSingle tests sorting single entry
func TestEntrySorterSingle(t *testing.T) {
	entries := EntrySorter{
		createMockDirEntry("single.md", false, time.Now()),
	}

	sort.Sort(entries)

	if len(entries) != 1 {
		t.Error("Single entry should remain single")
	}
	if entries[0].Name() != "single.md" {
		t.Errorf("Single entry name should be 'single.md', got '%s'", entries[0].Name())
	}
}

// mockDirEntry implements fs.DirEntry for testing
type mockDirEntry struct {
	name    string
	isDir   bool
	modTime time.Time
}

func createMockDirEntry(name string, isDir bool, modTime time.Time) *mockDirEntry {
	return &mockDirEntry{name: name, isDir: isDir, modTime: modTime}
}

func (m *mockDirEntry) Name() string { return m.name }
func (m *mockDirEntry) IsDir() bool  { return m.isDir }
func (m *mockDirEntry) Type() fs.FileMode {
	if m.isDir {
		return fs.ModeDir
	}
	return 0
}
func (m *mockDirEntry) Info() (fs.FileInfo, error) {
	return &mockFileInfo{name: m.name, isDir: m.isDir, modTime: m.modTime}, nil
}

// mockFileInfo implements fs.FileInfo for testing
type mockFileInfo struct {
	name    string
	isDir   bool
	modTime time.Time
}

func (m *mockFileInfo) Name() string { return m.name }
func (m *mockFileInfo) Size() int64  { return 0 }
func (m *mockFileInfo) Mode() fs.FileMode {
	if m.isDir {
		return fs.ModeDir
	}
	return 0
}
func (m *mockFileInfo) ModTime() time.Time { return m.modTime }
func (m *mockFileInfo) IsDir() bool        { return m.isDir }
func (m *mockFileInfo) Sys() interface{}   { return nil }
