package serverHandler

import (
	"bytes"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"mjpclab.dev/ghfs/src/serverLog"
)

func TestLogEscapesControllingRune(t *testing.T) {
	buf := &bytes.Buffer{}
	wMan := serverLog.NewWriterMan()
	h := &aliasHandler{logger: wMan.NewLogger(buf, buf, nil)}
	r := httptest.NewRequest("POST", "/", nil)

	h.logMutate("alice", "delete", "/tmp/a\nforged", r)
	h.logUpload("alice", "b\r\nforged", "/tmp/b\nforged", r)
	h.logArchive("c\nforged.zip", "d\nforged", r)
	wMan.Close()

	lines := strings.Split(strings.TrimSuffix(buf.String(), "\n"), "\n")
	if len(lines) != 3 {
		t.Fatalf("expect 3 lines, got %d:\n%s", len(lines), buf.String())
	}
	for _, line := range lines {
		if !strings.Contains(line, `\nforged`) {
			t.Error(line)
		}
	}
}

func TestLogMutationToMutationLog(t *testing.T) {
	dir := t.TempDir()
	accPath := filepath.Join(dir, "access.log")
	mutPath := filepath.Join(dir, "mutation.log")

	man := serverLog.NewHybridMan()
	logger, es := man.NewLogger(accPath, mutPath, "")
	if len(es) > 0 {
		t.Fatal(es)
	}
	h := &aliasHandler{logger: logger}
	r := httptest.NewRequest("POST", "/", nil)

	h.logMutate("alice", "mkdir", "/tmp/a", r)
	h.logMutate("alice", "delete", "/tmp/b", r)
	h.logUpload("alice", "c", "/tmp/c", r)
	h.logArchive("d.zip", "d", r)
	man.Close()

	mutLog, _ := os.ReadFile(mutPath)
	mutLines := strings.Split(strings.TrimSuffix(string(mutLog), "\n"), "\n")
	if len(mutLines) != 3 {
		t.Errorf("expect 3 mutation lines, got %d:\n%s", len(mutLines), mutLog)
	}
	accLog, _ := os.ReadFile(accPath)
	if !strings.Contains(string(accLog), " archive file: ") || strings.Count(string(accLog), "\n") != 1 {
		t.Errorf("expect only archive line in access log:\n%s", accLog)
	}
}
