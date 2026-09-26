package serverHandler

import (
	"bytes"
	"net/http/httptest"
	"strings"
	"testing"

	"mjpclab.dev/ghfs/src/serverLog"
)

func TestLogEscapesControllingRune(t *testing.T) {
	buf := &bytes.Buffer{}
	wMan := serverLog.NewWriterMan()
	h := &aliasHandler{logger: wMan.NewLogger(buf, nil)}
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
