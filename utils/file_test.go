package utils

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestCacheDir_SeparatesOutputDirs(t *testing.T) {
	a := CacheDir(filepath.Join(t.TempDir(), "a", "output"))
	b := CacheDir(filepath.Join(t.TempDir(), "b", "output"))
	if a == b {
		t.Fatalf("different output dirs share cache dir %s", a)
	}
	out := filepath.Join(t.TempDir(), "output")
	if CacheDir(out) != CacheDir(out+string(filepath.Separator)) {
		t.Error("same output dir mapped to different cache dirs")
	}
	if filepath.Base(filepath.Dir(a)) != "subs-check" {
		t.Errorf("unexpected cache dir layout: %s", a)
	}
}

func TestWriteFileAtomic_CreatesDirAndReplaces(t *testing.T) {
	path := filepath.Join(t.TempDir(), "a", "b", "data.json")

	if err := WriteFileAtomic(path, []byte("one")); err != nil {
		t.Fatalf("first write: %v", err)
	}
	if err := WriteFileAtomic(path, []byte("two")); err != nil {
		t.Fatalf("second write: %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil || string(got) != "two" {
		t.Fatalf("read = %q, %v; want two", got, err)
	}

	entries, err := os.ReadDir(filepath.Dir(path))
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Fatalf("leftover temp files: %v", entries)
	}

	if runtime.GOOS != "windows" {
		fi, err := os.Stat(path)
		if err != nil {
			t.Fatal(err)
		}
		if perm := fi.Mode().Perm(); perm != 0o600 {
			t.Errorf("perm = %o, want 600", perm)
		}
	}
}
