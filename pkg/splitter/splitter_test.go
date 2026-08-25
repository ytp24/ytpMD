package splitter

import (
	"strings"
	"testing"

	"github.com/ytp24/ytp24/pkg/core"
)

func TestSplitter_SplitIntoChapters(t *testing.T) {
	cfg := core.DefaultConfig()
	s := NewSplitter(cfg)

	pages := []core.PDFPage{
		{
			PageNumber: 1,
			Lines: []string{
				"CHAPTER 1: INTRODUCTION TO CONTAINERS",
				"Containers isolate applications at the OS level.",
			},
		},
		{
			PageNumber: 2,
			Lines: []string{
				"CHAPTER 2: KUBERNETES ARCHITECTURE",
				"Kubernetes schedules containers across a cluster.",
			},
		},
	}

	chapters := s.SplitIntoChapters(pages, "TestBook")

	if len(chapters) != 2 {
		t.Fatalf("expected 2 chapters, got %d", len(chapters))
	}

	if chapters[0].Filename != "01_introduction_to_containers.md" {
		t.Errorf("expected filename '01_introduction_to_containers.md', got '%s'", chapters[0].Filename)
	}

	if chapters[1].Filename != "02_kubernetes_architecture.md" {
		t.Errorf("expected filename '02_kubernetes_architecture.md', got '%s'", chapters[1].Filename)
	}

	toc := s.GenerateTOCIndex("TestBook", chapters, 2)
	if !strings.Contains(toc, "01_introduction_to_containers.md") {
		t.Errorf("expected TOC to contain link to chapter 1")
	}
}
