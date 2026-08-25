package splitter

import (
	"testing"

	"github.com/ytp24/ytpMD/pkg/core"
)

func TestSplitter_TOCPageIgnore(t *testing.T) {
	cfg := core.DefaultConfig()
	s := NewSplitter(cfg)

	pages := []core.PDFPage{
		{
			PageNumber: 1,
			Lines: []string{
				"Table of Contents",
				"Chapter 1: Cloud Fundamentals .................. 1",
				"Chapter 2: Kubernetes Orchestration ............ 25",
				"Chapter 3: AWS Risk and Compliance ............. 50",
			},
		},
		{
			PageNumber: 2,
			Lines: []string{
				"CHAPTER 1: CLOUD FUNDAMENTALS",
				"Cloud computing delivers computing services over the internet.",
				"It provides scalable and on-demand compute infrastructure for modern organizations.",
			},
		},
		{
			PageNumber: 3,
			Lines: []string{
				"CHAPTER 2: KUBERNETES ORCHESTRATION",
				"Kubernetes automates container lifecycle management across clusters.",
				"It manages deployment, scaling, service discovery, and zero-downtime rollouts.",
			},
		},
	}

	chapters := s.SplitIntoChapters(pages, "StudyGuide")

	// Must extract 2 actual chapters, NOT 5 (TOC lines ignored)
	if len(chapters) != 2 {
		t.Fatalf("expected 2 actual body chapters, got %d", len(chapters))
	}

	if chapters[0].Title != "CHAPTER 1: CLOUD FUNDAMENTALS" {
		t.Errorf("expected Chapter 1 title, got '%s'", chapters[0].Title)
	}

	if chapters[1].Title != "CHAPTER 2: KUBERNETES ORCHESTRATION" {
		t.Errorf("expected Chapter 2 title, got '%s'", chapters[1].Title)
	}
}
