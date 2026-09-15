package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
)

// CacheDir returns the private cache dir (results snapshot, export cache) of the
// instance writing to outputDir. Exports contain credentials, so it lives in the
// per-user cache dir rather than the shared /tmp or the public output dir, keyed
// by the output path so separate instances never share it.
func CacheDir(outputDir string) string {
	base, err := os.UserCacheDir()
	if err != nil {
		base = os.TempDir()
	}
	if abs, err := filepath.Abs(outputDir); err == nil {
		outputDir = abs
	}
	sum := sha256.Sum256([]byte(outputDir))
	return filepath.Join(base, "subs-check", hex.EncodeToString(sum[:6]))
}

// WriteFileAtomic writes via a temp file and rename so readers never see partial
// content. Missing dirs are created with 0700; the file ends up 0600.
func WriteFileAtomic(path string, data []byte) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(dir, ".tmp-"+filepath.Base(path)+"-*")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		os.Remove(tmpName)
		return err
	}
	if err := tmp.Close(); err != nil {
		os.Remove(tmpName)
		return err
	}
	if err := os.Rename(tmpName, path); err != nil {
		os.Remove(tmpName)
		return err
	}
	return nil
}
