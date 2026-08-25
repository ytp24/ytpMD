package transformer

import (
	"strings"
	"testing"

	"github.com/ytp24/ytpMD/pkg/core"
)

func TestTransformer_NoCodeBlockOnIndentedHeadings(t *testing.T) {
	cfg := core.DefaultConfig()
	tr := NewTransformer(cfg)

	input := [][]string{
		{
			"    CHAPTER 13: AWS RISK AND COMPLIANCE",
			"    This chapter explains risk management and compliance programs on AWS.",
			"    Organizations must ensure cloud governance across multi-account structures.",
			"",
			"    13.1 Shared Responsibility Model",
			"    • AWS manages security of the cloud",
			"    • Customer manages security in the cloud",
			"",
			"    $ kubectl get pods -n kube-system",
			"    $ aws s3 ls",
		},
	}

	result := tr.Transform(input)

	// Heading must NOT be wrapped in ```bash
	if strings.Contains(result, "```bash\n# CHAPTER 13") || strings.Contains(result, "```bash\n\n# CHAPTER 13") {
		t.Fatalf("heading was improperly wrapped in a ```bash code block:\n%s", result)
	}

	if !strings.Contains(result, "# CHAPTER 13: AWS RISK AND COMPLIANCE") {
		t.Fatalf("expected Markdown heading '# CHAPTER 13...', got:\n%s", result)
	}

	// Real commands SHOULD be detected in code blocks
	if !strings.Contains(result, "```bash") || !strings.Contains(result, "kubectl get pods") {
		t.Fatalf("expected code block for shell commands, got:\n%s", result)
	}
}
