# 🤖 Model Context Protocol (MCP) Integration Guide

`ytpMD` embeds a native **Model Context Protocol (MCP)** server operating over `stdio` JSON-RPC 2.0. This allows AI coding agents and LLM-powered IDEs (Cursor, Claude Desktop, Google Antigravity, VS Code, Zed, Windsurf) to autonomously convert and index technical PDF documentation into chapter-split Markdown context.

---

## ⚡ Quick Start

Launch the MCP server in stdio mode:
```bash
ytpmd mcp
```
*(Or invoke via the dedicated binary `ytpmd-mcp`)*

---

## 🛠️ Registered MCP Tools

| Tool | Parameters | Description |
| :--- | :--- | :--- |
| `convert_pdf` | `path` (req), `output_dir`, `single_file`, `skip_front_matter`, `start_page`, `end_page`, `exclude_appendix`, `reflow_paragraphs` | Converts a PDF into a chapter-split Markdown folder with YAML frontmatter & `AGENTS.md` manifest. |
| `batch_convert` | `directory` (req), `output_dir`, `batch_name`, `concurrency`, `recursive`, `exclude_appendix` | Concurrently processes an entire folder of PDFs using multi-core worker pool. |
| `inspect_pdf` | `path` (req) | Validates `%PDF-` header, page count, encryption status, and text layer readability without disk writes. |
| `get_manifest` | `folder_path` (req) | Retrieves the machine-readable `AGENTS.md` context map and chapter index for an extracted book. |

---

## ⚙️ IDE Configuration Snippets

### 1. Cursor (`.cursor/mcp.json`)
```json
{
  "mcpServers": {
    "ytpmd": {
      "command": "ytpmd",
      "args": ["mcp"]
    }
  }
}
```

### 2. Claude Desktop (`claude_desktop_config.json`)
- **macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`
- **Windows**: `%APPDATA%\Claude\claude_desktop_config.json`
- **Linux**: `~/.config/Claude/claude_desktop_config.json`

```json
{
  "mcpServers": {
    "ytpmd": {
      "command": "ytpmd",
      "args": ["mcp"]
    }
  }
}
```

### 3. Google Antigravity (`~/.gemini/antigravity-cli/settings.json`)
```json
{
  "mcp": {
    "servers": {
      "ytpmd": {
        "command": "ytpmd",
        "args": ["mcp"]
      }
    }
  }
}
```

### 4. VS Code (Continue / Roo Code) (`.vscode/settings.json`)
```json
{
  "mcp.servers": {
    "ytpmd": {
      "command": "ytpmd",
      "args": ["mcp"]
    }
  }
}
```

### 5. Zed Editor (`~/.config/zed/settings.json`)
```json
{
  "context_servers": {
    "ytpmd": {
      "command": "ytpmd",
      "args": ["mcp"]
    }
  }
}
```

---

## 🔒 Security & Privacy
- Runs 100% offline via local process pipes (`stdin` / `stdout`).
- Zero external network telemetry.
- Defensively checks file headers and enforces strict execution timeouts on child processes.
