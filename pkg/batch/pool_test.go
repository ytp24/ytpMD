package batch

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/devops/pdf2md/pkg/core"
)

// MockExtractor simulates PDF extraction for unit tests.
type MockExtractor struct{}

func (m *MockExtractor) ExtractFile(ctx context.Context, path string) (*core.ProcessedDocument, error) {
	return &core.ProcessedDocument{
		SourcePath:      path,
		TotalPages:      4,
		ProcessedPages:  3,
		MarkdownContent: "# Mock Content",
	}, nil
}

func (m *MockExtractor) ExtractToDirectory(ctx context.Context, path string, targetBaseDir string) (*core.SplitResult, error) {
	base := filepath.Base(path)
	if base == "corrupt.pdf" {
		return nil, fmt.Errorf("corrupt file")
	}

	targetDir := filepath.Join(targetBaseDir, base)
	_ = os.MkdirAll(targetDir, 0755)

	return &core.SplitResult{
		SourcePDF:       path,
		PDFName:         base,
		TargetDirectory: targetDir,
		TotalPages:      4,
		ProcessedPages:  3,
		Chapters: []core.Chapter{
			{Index: 1, Title: "Intro", Filename: "01_intro.md"},
			{Index: 2, Title: "Details", Filename: "02_details.md"},
		},
	}, nil
}

func TestConcurrentBatchEngine_ProcessBatch(t *testing.T) {
	tempDir := t.TempDir()
	mockExt := &MockExtractor{}
	engine := NewConcurrentBatchEngine(mockExt)

	pdfFiles := []string{
		filepath.Join(tempDir, "doc1.pdf"),
		filepath.Join(tempDir, "doc2.pdf"),
		filepath.Join(tempDir, "doc3.pdf"),
		filepath.Join(tempDir, "corrupt.pdf"),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := engine.ProcessBatch(ctx, pdfFiles, tempDir, "test_batch", 2, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.TotalFiles != 4 {
		t.Errorf("expected 4 total files, got %d", result.TotalFiles)
	}

	if result.ProcessedFiles != 3 {
		t.Errorf("expected 3 processed files, got %d", result.ProcessedFiles)
	}

	if result.FailedFiles != 1 {
		t.Errorf("expected 1 failed file, got %d", result.FailedFiles)
	}

	// Verify master README.md was generated
	masterReadme := filepath.Join(tempDir, "test_batch", "README.md")
	if _, err := os.Stat(masterReadme); os.IsNotExist(err) {
		t.Errorf("expected master README.md to be created at %s", masterReadme)
	}
}
