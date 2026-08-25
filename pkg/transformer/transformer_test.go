package transformer

import (
	"strings"
	"testing"

	"github.com/devops/pdf2md/pkg/models"
)

func TestTransformer_Transform(t *testing.T) {
	cfg := models.DefaultConfig()
	tr := NewTransformer(cfg)

	input := [][]string{
		{
			"CHAPTER 1: CLOUD ARCHITECTURE",
			"This is a high-perfor-",
			"mance cloud computing archi-",
			"tecture for enterprise systems.",
			"",
			"1.1 Control Plane",
			"• First bullet point",
			"▪ Second bullet point",
		},
	}

	result := tr.Transform(input)

	if !strings.Contains(result, "# CHAPTER 1: CLOUD ARCHITECTURE") {
		t.Fatalf("expected H1 Chapter header, got:\n%s", result)
	}

	if !strings.Contains(result, "high-performance cloud computing architecture for enterprise systems.") {
		t.Fatalf("expected reflowed and de-hyphenated paragraph, got:\n%s", result)
	}

	if !strings.Contains(result, "## 1.1 Control Plane") {
		t.Fatalf("expected H2 Section header, got:\n%s", result)
	}

	if !strings.Contains(result, "- First bullet point") || !strings.Contains(result, "- Second bullet point") {
		t.Fatalf("expected standardized bullet points, got:\n%s", result)
	}
}
