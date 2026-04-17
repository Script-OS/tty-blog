package global

import (
	"testing"
)

// TestDefaultConfig tests that defaultConfig returns valid defaults
func TestDefaultConfig(t *testing.T) {
	cfg := defaultConfig()

	if cfg.Editor == nil {
		t.Error("Editor should not be nil in default config")
	}

	if len(*cfg.Editor) == 0 {
		t.Error("Editor should have at least one element")
	}

	if cfg.RootDir == nil {
		t.Error("RootDir should not be nil in default config")
	}

	if *cfg.RootDir != "." {
		t.Errorf("RootDir should be '.', got %s", *cfg.RootDir)
	}
}

// TestMergeConfig tests configuration merging logic
func TestMergeConfig(t *testing.T) {
	editor := []string{"vim"}
	password := "secret"
	rootDir := "/home/user/blog"

	userCfg := &ConfigType{
		Editor:         &editor,
		EditorPassword: &password,
		RootDir:        &rootDir,
	}

	defaultCfg := defaultConfig()

	// Merge with user config first (higher priority)
	result := mergeConfig(userCfg, defaultCfg)

	if result.Editor == nil || (*result.Editor)[0] != "vim" {
		t.Errorf("Editor should be vim, got %v", result.Editor)
	}

	if result.EditorPassword == nil || *result.EditorPassword != "secret" {
		t.Errorf("EditorPassword should be secret, got %v", result.EditorPassword)
	}

	if result.RootDir == nil || *result.RootDir != "/home/user/blog" {
		t.Errorf("RootDir should be /home/user/blog, got %v", result.RootDir)
	}
}

// TestMergeConfigPartial tests partial configuration merging
func TestMergeConfigPartial(t *testing.T) {
	editor := []string{"code", "--wait"}

	// User config only specifies editor
	userCfg := &ConfigType{
		Editor: &editor,
	}

	defaultCfg := defaultConfig()

	result := mergeConfig(userCfg, defaultCfg)

	if result.Editor == nil || len(*result.Editor) != 2 {
		t.Errorf("Editor should have 2 elements, got %v", result.Editor)
	}

	// RootDir should come from default
	if result.RootDir == nil || *result.RootDir != "." {
		t.Errorf("RootDir should be '.' from default, got %v", result.RootDir)
	}
}

// TestMergeConfigEmpty tests merging with empty configs
func TestMergeConfigEmpty(t *testing.T) {
	emptyCfg := &ConfigType{}
	defaultCfg := defaultConfig()

	result := mergeConfig(emptyCfg, defaultCfg)

	// All values should come from default
	if result.Editor == nil {
		t.Error("Editor should not be nil")
	}

	if result.RootDir == nil {
		t.Error("RootDir should not be nil")
	}
}
