package archive

import (
	"compress/gzip"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"archive/tar"
)

// TestExtract_ValidTarGz tests successful extraction of a valid .tar.gz file.
func TestExtract_ValidTarGz(t *testing.T) {
	// Create temporary directories
	srcDir := t.TempDir()
	destDir := t.TempDir()

	tarGzPath := filepath.Join(srcDir, "test.tar.gz")

	// Create a valid .tar.gz with some files
	if err := createTestTarGz(tarGzPath, map[string]string{
		"file.txt":    "hello world",
		"subdir/data": "test data",
	}); err != nil {
		t.Fatalf("failed to create test archive: %v", err)
	}

	// Extract
	if err := Extract(tarGzPath, destDir); err != nil {
		t.Fatalf("Extract() failed: %v", err)
	}

	// Verify files exist and have correct content
	testFile := filepath.Join(destDir, "file.txt")
	if content, err := os.ReadFile(testFile); err != nil || string(content) != "hello world" {
		t.Errorf("file.txt not extracted correctly")
	}

	dataFile := filepath.Join(destDir, "subdir", "data")
	if content, err := os.ReadFile(dataFile); err != nil || string(content) != "test data" {
		t.Errorf("subdir/data not extracted correctly")
	}
}

// TestExtract_InvalidMagicBytes tests rejection of files with invalid gzip magic bytes.
func TestExtract_InvalidMagicBytes(t *testing.T) {
	srcDir := t.TempDir()
	destDir := t.TempDir()

	// Create a file with invalid magic bytes
	badFile := filepath.Join(srcDir, "notgzip.tar.gz")
	if err := os.WriteFile(badFile, []byte("this is not a gzip file"), 0644); err != nil {
		t.Fatalf("failed to create bad file: %v", err)
	}

	// Extract should fail with magic byte error
	err := Extract(badFile, destDir)
	if err == nil {
		t.Fatal("Extract() should have failed for non-gzip file")
	}
	if err.Error() != "not a valid gzip file (magic bytes mismatch)" {
		t.Errorf("unexpected error: %v", err)
	}
}

// TestExtract_DirectoryTraversal tests protection against directory traversal attacks.
func TestExtract_DirectoryTraversal(t *testing.T) {
	srcDir := t.TempDir()
	destDir := t.TempDir()

	tarGzPath := filepath.Join(srcDir, "malicious.tar.gz")

	// Create a .tar.gz with a directory traversal path
	if err := createTestTarGzWithMaliciousPath(tarGzPath, "../../../etc/passwd"); err != nil {
		t.Fatalf("failed to create test archive: %v", err)
	}

	// Extract should fail
	err := Extract(tarGzPath, destDir)
	if err == nil {
		t.Fatal("Extract() should have rejected directory traversal path")
	}
	if err.Error() != "unsafe path in archive: ../../../etc/passwd" {
		t.Errorf("unexpected error: %v", err)
	}
}

// TestExtract_MissingFile tests handling of non-existent archive files.
func TestExtract_MissingFile(t *testing.T) {
	destDir := t.TempDir()

	err := Extract("/nonexistent/path/archive.tar.gz", destDir)
	if err == nil {
		t.Fatal("Extract() should have failed for missing file")
	}
}

// TestExtract_DestDirCreation tests that destination directory is created if missing.
func TestExtract_DestDirCreation(t *testing.T) {
	srcDir := t.TempDir()
	destDir := filepath.Join(t.TempDir(), "nested", "dest")

	tarGzPath := filepath.Join(srcDir, "test.tar.gz")
	if err := createTestTarGz(tarGzPath, map[string]string{
		"file.txt": "content",
	}); err != nil {
		t.Fatalf("failed to create test archive: %v", err)
	}

	// Extract to non-existent nested directory
	if err := Extract(tarGzPath, destDir); err != nil {
		t.Fatalf("Extract() failed: %v", err)
	}

	// Verify destination was created
	if _, err := os.Stat(destDir); os.IsNotExist(err) {
		t.Fatal("destination directory was not created")
	}
}

// TestExtract_FilePermissions tests that file permissions are preserved.
func TestExtract_FilePermissions(t *testing.T) {
	srcDir := t.TempDir()
	destDir := t.TempDir()

	tarGzPath := filepath.Join(srcDir, "test.tar.gz")

	// Create .tar.gz with specific file permissions (0755)
	if err := createTestTarGzWithPermissions(tarGzPath, "executable.sh", "#!/bin/bash", 0755); err != nil {
		t.Fatalf("failed to create test archive: %v", err)
	}

	if err := Extract(tarGzPath, destDir); err != nil {
		t.Fatalf("Extract() failed: %v", err)
	}

	// Check file permissions
	execFile := filepath.Join(destDir, "executable.sh")
	info, err := os.Stat(execFile)
	if err != nil {
		t.Fatalf("failed to stat extracted file: %v", err)
	}

	if info.Mode()&0755 != 0755 {
		t.Errorf("file permissions not preserved: got %o, expected 0755", info.Mode())
	}
}

// TestValidateGzipMagic tests magic byte validation.
func TestValidateGzipMagic(t *testing.T) {
	tests := []struct {
		name      string
		magic     []byte
		wantError bool
	}{
		{"valid gzip", []byte{0x1f, 0x8b}, false},
		{"invalid magic", []byte{0xff, 0xff}, true},
		{"empty file", []byte{}, true},
		{"single byte", []byte{0x1f}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpFile := filepath.Join(t.TempDir(), "test")
			if err := os.WriteFile(tmpFile, tt.magic, 0644); err != nil {
				t.Fatalf("failed to create test file: %v", err)
			}

			err := validateGzipMagic(tmpFile)
			if (err != nil) != tt.wantError {
				t.Errorf("validateGzipMagic() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

// TestIsPathSafe tests directory traversal protection.
func TestIsPathSafe(t *testing.T) {
	destDir := t.TempDir()

	tests := []struct {
		name       string
		targetPath string
		wantSafe   bool
	}{
		{"safe: file in dest", filepath.Join(destDir, "file.txt"), true},
		{"safe: nested file", filepath.Join(destDir, "sub", "dir", "file.txt"), true},
		{"unsafe: parent traversal", filepath.Join(destDir, "..", "etc", "passwd"), false},
		{"unsafe: absolute path", "/etc/passwd", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			safe := isPathSafe(destDir, tt.targetPath)
			if safe != tt.wantSafe {
				t.Errorf("isPathSafe() = %v, want %v", safe, tt.wantSafe)
			}
		})
	}
}

// Helper functions

// createTestTarGz creates a valid .tar.gz file with the given files.
func createTestTarGz(path string, files map[string]string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Printf("warning: failed to close file: %v", err)
		}
	}()

	gzipWriter := gzip.NewWriter(file)
	defer func() {
		if err := gzipWriter.Close(); err != nil {
			fmt.Printf("warning: failed to close gzip writer: %v", err)
		}
	}()

	tarWriter := tar.NewWriter(gzipWriter)
	defer func() {
		if err := tarWriter.Close(); err != nil {
			fmt.Printf("warning: failed to close gzip writer: %v", err)
		}
	}()

	for name, content := range files {
		// Create parent directories if needed
		dir := filepath.Dir(name)
		if dir != "." {
			header := &tar.Header{
				Name:     dir,
				Typeflag: tar.TypeDir,
				Mode:     0755,
			}
			if err := tarWriter.WriteHeader(header); err != nil {
				return err
			}
		}

		// Add file
		header := &tar.Header{
			Name: name,
			Size: int64(len(content)),
			Mode: 0644,
		}
		if err := tarWriter.WriteHeader(header); err != nil {
			return err
		}
		if _, err := tarWriter.Write([]byte(content)); err != nil {
			return err
		}
	}

	return nil
}

// createTestTarGzWithMaliciousPath creates a .tar.gz with a path traversal entry.
func createTestTarGzWithMaliciousPath(path, maliciousPath string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Printf("warning: failed to close file: %v", err)
		}
	}()

	gzipWriter := gzip.NewWriter(file)
	defer func() {
		if err := gzipWriter.Close(); err != nil {
			fmt.Printf("warning: failed to close gzip writer: %v", err)
		}
	}()

	tarWriter := tar.NewWriter(gzipWriter)
	defer func() {
		if err := tarWriter.Close(); err != nil {
			fmt.Printf("warning: failed to close gzip writer: %v", err)
		}
	}()

	header := &tar.Header{
		Name: maliciousPath,
		Size: 4,
		Mode: 0644,
	}
	if err := tarWriter.WriteHeader(header); err != nil {
		return err
	}
	if _, err := tarWriter.Write([]byte("evil")); err != nil {
		return err
	}

	return tarWriter.Close()
}

// createTestTarGzWithPermissions creates a .tar.gz with a file having specific permissions.
func createTestTarGzWithPermissions(path, name, content string, mode int64) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Printf("warning: failed to close file: %v", err)
		}
	}()

	gzipWriter := gzip.NewWriter(file)
	defer func() {
		if err := gzipWriter.Close(); err != nil {
			fmt.Printf("warning: failed to close gzip writer: %v", err)
		}
	}()

	tarWriter := tar.NewWriter(gzipWriter)
	defer func() {
		if err := tarWriter.Close(); err != nil {
			fmt.Printf("warning: failed to close gzip writer: %v", err)
		}
	}()

	header := &tar.Header{
		Name: name,
		Size: int64(len(content)),
		Mode: mode,
	}
	if err := tarWriter.WriteHeader(header); err != nil {
		return err
	}
	if _, err := tarWriter.Write([]byte(content)); err != nil {
		return err
	}

	return tarWriter.Close()
}
