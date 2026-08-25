package validator

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidatePDFFile(t *testing.T) {
	tempDir := t.TempDir()

	// 1. Non-existent file
	err := ValidatePDFFile(filepath.Join(tempDir, "missing.pdf"))
	if err == nil {
		t.Errorf("expected error for non-existent file, got nil")
	}

	// 2. Empty file (0 bytes)
	emptyFile := filepath.Join(tempDir, "empty.pdf")
	_ = os.WriteFile(emptyFile, []byte{}, 0644)
	err = ValidatePDFFile(emptyFile)
	if err == nil {
		t.Errorf("expected error for 0-byte file, got nil")
	}

	// 3. Fake non-PDF file
	fakeFile := filepath.Join(tempDir, "fake.pdf")
	_ = os.WriteFile(fakeFile, []byte("This is plain text without PDF header"), 0644)
	err = ValidatePDFFile(fakeFile)
	if err == nil {
		t.Errorf("expected error for non-PDF file, got nil")
	}

	// 4. Valid mock PDF header
	validFile := filepath.Join(tempDir, "valid.pdf")
	_ = os.WriteFile(validFile, []byte("%PDF-1.4 mock content"), 0644)
	err = ValidatePDFFile(validFile)
	if err != nil {
		t.Errorf("expected valid PDF header to pass, got error: %v", err)
	}
}

func TestValidateDirectory(t *testing.T) {
	tempDir := t.TempDir()

	// 1. Existing writable directory
	err := ValidateDirectory(tempDir)
	if err != nil {
		t.Errorf("expected valid directory, got: %v", err)
	}

	// 2. Create nested directory
	nested := filepath.Join(tempDir, "nested", "subfolder")
	err = ValidateDirectory(nested)
	if err != nil {
		t.Errorf("expected nested directory creation to pass, got: %v", err)
	}
}
