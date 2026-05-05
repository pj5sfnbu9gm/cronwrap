package history_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourorg/cronwrap/internal/history"
)

func tempStore(t *testing.T) (*history.Store, string) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "history.json")
	s, err := history.NewStore(path)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	return s, path
}

func TestAppendAndRetrieve(t *testing.T) {
	s, _ := tempStore(t)

	rec := history.Record{
		JobName:    "backup",
		Command:    "tar -czf /tmp/backup.tar.gz /data",
		StartedAt:  time.Now().UTC(),
		FinishedAt: time.Now().UTC(),
		Duration:   2 * time.Second,
		ExitCode:   0,
		Success:    true,
	}
	if err := s.Append(rec); err != nil {
		t.Fatalf("Append: %v", err)
	}

	records := s.Records()
	if len(records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(records))
	}
	if records[0].JobName != "backup" {
		t.Errorf("expected job_name=backup, got %s", records[0].JobName)
	}
}

func TestPersistence(t *testing.T) {
	s, path := tempStore(t)

	for i := 0; i < 3; i++ {
		_ = s.Append(history.Record{JobName: "job", Success: true})
	}

	s2, err := history.NewStore(path)
	if err != nil {
		t.Fatalf("reopen store: %v", err)
	}
	if len(s2.Records()) != 3 {
		t.Errorf("expected 3 records after reload, got %d", len(s2.Records()))
	}
}

func TestNewStoreCreatesDir(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sub", "dir", "history.json")
	s, err := history.NewStore(path)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	_ = s.Append(history.Record{JobName: "test"})
	if _, err := os.Stat(path); err != nil {
		t.Errorf("expected file to exist: %v", err)
	}
}
