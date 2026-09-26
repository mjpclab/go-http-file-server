package serverHandler

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestGetAvailableFilename(t *testing.T) {
	dir := t.TempDir()

	check := func(filename string, mustAppendSuffix bool, expected string) {
		t.Helper()
		result, err := getAvailableFilename(dir, filename, mustAppendSuffix)
		if err != nil || result != expected {
			t.Errorf("%s: got %q %v, expect %q", filename, result, err, expected)
		}
	}

	check("a.txt", false, "a.txt")
	check("a.txt", true, "a-1.txt")

	if err := os.WriteFile(filepath.Join(dir, "a.txt"), nil, 0644); err != nil {
		t.Fatal(err)
	}
	check("a.txt", false, "a-1.txt")

	if err := os.WriteFile(filepath.Join(dir, "a-1.txt"), nil, 0644); err != nil {
		t.Fatal(err)
	}
	check("a.txt", false, "a-2.txt")
	check("a.txt", true, "a-2.txt")
}

func TestGetAvailableFilenameError(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("ENOTDIR is not reported on windows")
	}

	notDir := filepath.Join(t.TempDir(), "file")
	if err := os.WriteFile(notDir, nil, 0644); err != nil {
		t.Fatal(err)
	}

	for _, mustAppendSuffix := range []bool{false, true} {
		result, err := getAvailableFilename(notDir, "a.txt", mustAppendSuffix)
		if err == nil || result != "" {
			t.Errorf("mustAppendSuffix=%v: got %q %v, expect error", mustAppendSuffix, result, err)
		}
	}
}
