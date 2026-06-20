package archive

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Extract decompresses and extracts a .tar.gz file to the destination directory.
// It validates the file format using magic bytes before proceeding.
func Extract(tarGzPath, destDir string) error {
	// Validate magic bytes
	if err := validateGzipMagic(tarGzPath); err != nil {
		return err
	}

	// Create destination directory if it doesn't exist
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Open the .tar.gz file
	file, err := os.Open(tarGzPath)
	if err != nil {
		return fmt.Errorf("failed to open archive: %w", err)
	}
	defer file.Close()

	// Decompress gzip
	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("invalid gzip format: %w", err)
	}
	defer gzipReader.Close()

	// Extract tar entries
	tarReader := tar.NewReader(gzipReader)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			return fmt.Errorf("failed to read tar entry: %w", err)
		}

		// Construct the full path and prevent directory traversal attacks
		targetPath := filepath.Join(destDir, header.Name)
		if !isPathSafe(destDir, targetPath) {
			return fmt.Errorf("unsafe path in archive: %s", header.Name)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(targetPath, header.FileInfo().Mode()); err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}
		case tar.TypeReg:
			// Create parent directory if needed
			if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}

			// Extract file
			outFile, err := os.Create(targetPath)
			if err != nil {
				return fmt.Errorf("failed to create file: %w", err)
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close()
				return fmt.Errorf("failed to extract file: %w", err)
			}
			outFile.Close()

			// Preserve file permissions
			if err := os.Chmod(targetPath, header.FileInfo().Mode()); err != nil {
				return fmt.Errorf("failed to set file permissions: %w", err)
			}
		}
	}

	return nil
}

// validateGzipMagic checks for gzip magic bytes (0x1f 0x8b) at the start of the file.
func validateGzipMagic(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		// File not found is skippable
		if errors.Is(err, os.ErrNotExist) {
			return NewSkippableError("archive not found")
		}
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Read first two bytes
	magic := make([]byte, 2)
	if _, err := file.Read(magic); err != nil {
		return fmt.Errorf("failed to read file header: %w", err)
	}

	// Check gzip magic bytes
	if magic[0] != 0x1f || magic[1] != 0x8b {
		return NewSkippableError("not a valid gzip file (magic bytes mismatch)")
	}

	return nil
}

// isPathSafe ensures that the extracted path stays within the destination directory
// to prevent directory traversal attacks (e.g., ../../../etc/passwd).
func isPathSafe(destDir, targetPath string) bool {
	// Resolve to absolute paths
	absDestDir, err := filepath.Abs(destDir)
	if err != nil {
		return false
	}
	absTargetPath, err := filepath.Abs(targetPath)
	if err != nil {
		return false
	}

	// Get relative path from destination to target
	relPath, err := filepath.Rel(absDestDir, absTargetPath)
	if err != nil {
		return false
	}

	// If relative path is absolute or starts with "..", it's outside the destination
	return !filepath.IsAbs(relPath) && !strings.HasPrefix(relPath, "..")
}
