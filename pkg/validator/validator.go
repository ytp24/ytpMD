package validator

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
)

// ExpandPath expands leading ~ to the user's home directory and returns an absolute clean path.
func ExpandPath(path string) string {
	path = strings.TrimSpace(path)
	if path == "" {
		return ""
	}

	if strings.HasPrefix(path, "~") {
		usr, err := user.Current()
		if err == nil {
			if path == "~" {
				path = usr.HomeDir
			} else if strings.HasPrefix(path, "~/") {
				path = filepath.Join(usr.HomeDir, path[2:])
			}
		}
	}

	abs, err := filepath.Abs(path)
	if err != nil {
		return filepath.Clean(path)
	}
	return abs
}

// CheckDependencies verifies that system utilities (pdftotext) are installed.
func CheckDependencies() error {
	_, err := exec.LookPath("pdftotext")
	if err != nil {
		return fmt.Errorf("missing system dependency 'pdftotext'.\nPlease install poppler-utils (e.g. 'sudo apt install poppler-utils' or 'brew install poppler')")
	}
	return nil
}

// ValidatePDFFile verifies the file exists, is readable, and contains PDF header magic bytes.
func ValidatePDFFile(path string) error {
	cleanPath := ExpandPath(path)

	info, err := os.Stat(cleanPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file does not exist: '%s'", path)
		}
		return fmt.Errorf("cannot access file: %w", err)
	}

	if info.IsDir() {
		return fmt.Errorf("path is a directory, not a PDF file: '%s'", path)
	}

	if info.Size() == 0 {
		return fmt.Errorf("file is empty (0 bytes): '%s'", path)
	}

	// Verify PDF magic bytes (%PDF-)
	file, err := os.Open(cleanPath)
	if err != nil {
		return fmt.Errorf("permission denied reading file: %w", err)
	}
	defer file.Close()

	header := make([]byte, 1024)
	n, err := file.Read(header)
	if err != nil || n < 4 {
		return fmt.Errorf("unable to read file header: %w", err)
	}

	if !bytes.Contains(header[:n], []byte("%PDF-")) {
		return fmt.Errorf("invalid file format: '%s' is not a valid PDF document (missing %%PDF- header)", filepath.Base(path))
	}

	return nil
}

// ValidateDirectory ensures the directory path is valid or can be created.
func ValidateDirectory(path string) error {
	cleanPath := ExpandPath(path)

	info, err := os.Stat(cleanPath)
	if err == nil {
		if !info.IsDir() {
			return fmt.Errorf("destination '%s' exists and is a file, not a directory", path)
		}
		// Test write permission
		testFile := filepath.Join(cleanPath, ".pdf2md_write_test")
		if err := os.WriteFile(testFile, []byte(""), 0644); err != nil {
			return fmt.Errorf("destination directory is not writable (permission denied): %w", err)
		}
		_ = os.Remove(testFile)
		return nil
	}

	if os.IsNotExist(err) {
		// Attempt creation
		if err := os.MkdirAll(cleanPath, 0755); err != nil {
			return fmt.Errorf("failed to create destination directory '%s': %w", path, err)
		}
		return nil
	}

	return fmt.Errorf("invalid directory path: %w", err)
}
