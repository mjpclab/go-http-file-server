package serverLog

import (
	"os"
	"path/filepath"
	"testing"
)

func TestHybridManMutationLog(t *testing.T) {
	dir := t.TempDir()
	accPath := filepath.Join(dir, "access.log")
	mutPath := filepath.Join(dir, "mutation.log")

	man := NewHybridMan()
	logger, es := man.NewLogger(accPath, mutPath, "")
	if len(es) > 0 {
		t.Fatal(es)
	}
	logger.LogAccessString("access")
	logger.LogMutationString("mutation")
	man.Close()

	accLog, _ := os.ReadFile(accPath)
	if string(accLog) != "access\n" {
		t.Error(string(accLog))
	}
	mutLog, _ := os.ReadFile(mutPath)
	if string(mutLog) != "mutation\n" {
		t.Error(string(mutLog))
	}
}

func TestHybridManMutationLogDisabled(t *testing.T) {
	accPath := filepath.Join(t.TempDir(), "access.log")

	man := NewHybridMan()
	logger, es := man.NewLogger(accPath, "", "")
	if len(es) > 0 {
		t.Fatal(es)
	}
	if logger.CanLogMutation() {
		t.Error("empty mutation log should be disabled")
	}
	logger.LogAccessString("access")
	logger.LogMutationString("mutation")
	man.Close()

	accLog, _ := os.ReadFile(accPath)
	if string(accLog) != "access\n" {
		t.Error(string(accLog))
	}
}

func TestHybridManMutationLogOnly(t *testing.T) {
	mutPath := filepath.Join(t.TempDir(), "mutation.log")

	man := NewHybridMan()
	logger, es := man.NewLogger("", mutPath, "")
	if len(es) > 0 {
		t.Fatal(es)
	}
	if logger.CanLogAccess() || !logger.CanLogMutation() {
		t.Error("logger should only log mutation")
	}
	logger.LogMutationString("mutation")
	man.Close()

	mutLog, _ := os.ReadFile(mutPath)
	if string(mutLog) != "mutation\n" {
		t.Error(string(mutLog))
	}
}
