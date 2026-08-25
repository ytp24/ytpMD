package filter

import (
	"strings"
	"testing"

	"github.com/ytp24/ytp24/pkg/core"
)

func TestFilter_ShouldStop(t *testing.T) {
	cfg := core.DefaultConfig()
	f := NewContentFilter(cfg)

	// Appendix test
	stop, reason := f.ShouldStop("APPENDIX A: SYSTEM METRICS\nSome details...")
	if !stop {
		t.Fatalf("expected stop on Appendix A, got false")
	}
	if !strings.Contains(strings.ToLower(reason), "appendix") {
		t.Fatalf("expected reason to contain 'appendix', got %s", reason)
	}

	// Index test
	stop, reason = f.ShouldStop("INDEX\nAlphabetical...")
	if !stop {
		t.Fatalf("expected stop on Index, got false")
	}

	// Normal chapter test
	stop, _ = f.ShouldStop("CHAPTER 1: INTRODUCTION\nWelcome...")
	if stop {
		t.Fatalf("expected no stop on Chapter 1, got true")
	}
}

func TestFilter_CleanPageLines(t *testing.T) {
	cfg := core.DefaultConfig()
	f := NewContentFilter(cfg)

	raw := "Page 12 of 45\nCopyright 2026 Corporation\n\nSome important content [image: logo.png]\nFigure 1.2: Diagram\n99"
	lines := f.CleanPageLines(raw)

	joined := strings.Join(lines, "\n")
	if strings.Contains(joined, "Page 12 of 45") {
		t.Fatalf("expected header 'Page 12 of 45' to be removed")
	}
	if strings.Contains(joined, "[image: logo.png]") {
		t.Fatalf("expected image placeholder to be removed")
	}
	if !strings.Contains(joined, "Some important content") {
		t.Fatalf("expected 'Some important content' to be preserved")
	}
}
