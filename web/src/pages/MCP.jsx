import React, { useState } from 'react';
import CodeBlock from '../components/CodeBlock';

export default function MCP() {
  const [ideTab, setIdeTab] = useState('cursor');

  const configs = {
    cursor: `// .cursor/mcp.json
{
  "mcpServers": {
    "ytpmd": {
      "command": "ytpmd",
      "args": ["mcp"]
    }
  }
}`,
    claude: `// claude_desktop_config.json
{
  "mcpServers": {
    "ytpmd": {
      "command": "ytpmd",
      "args": ["mcp"]
    }
  }
}`,
    antigravity: `// ~/.gemini/antigravity-cli/settings.json
{
  "mcp": {
    "servers": {
      "ytpmd": {
        "command": "ytpmd",
        "args": ["mcp"]
      }
    }
  }
}`,
    vscode: `// .vscode/settings.json (Continue / Roo Code)
{
  "mcp.servers": {
    "ytpmd": {
      "command": "ytpmd",
      "args": ["mcp"]
    }
  }
}`,
    zed: `// ~/.config/zed/settings.json
{
  "context_servers": {
    "ytpmd": {
      "command": "ytpmd",
      "args": ["mcp"]
    }
  }
}`
  };

  return (
    <div>
      <div className="pixel-box">
        <div className="pixel-box-header">
          <span>/// 04_MCP_SPECIFICATION</span>
          <span>&gt; JSON-RPC_2.0</span>
        </div>

        <div className="pixel-box-body">
          <p style={{ color: '#cbd5e1', marginBottom: '14px' }}>
            ytpMD exposes a high-throughput Model Context Protocol (MCP) server over standard input/output (<code>stdio</code>). AI agents can autonomously inspect, convert, and index books in real-time.
          </p>

          <div className="pixel-section-header">
            <span>REGISTERED MCP TOOLS</span>
          </div>

          <table className="pixel-table" style={{ marginBottom: '24px' }}>
            <thead>
              <tr>
                <th>Tool Name</th>
                <th>Input Schema</th>
                <th>Output</th>
              </tr>
            </thead>
            <tbody>
              <tr>
                <td><code>convert_pdf</code></td>
                <td><code>path</code>, <code>output_dir</code>, <code>single_file</code>, <code>skip_front_matter</code></td>
                <td>Chapter index & token breakdown.</td>
              </tr>
              <tr>
                <td><code>batch_convert</code></td>
                <td><code>directory</code>, <code>output_dir</code>, <code>concurrency</code>, <code>recursive</code></td>
                <td>Concurrent batch extraction metrics.</td>
              </tr>
              <tr>
                <td><code>inspect_pdf</code></td>
                <td><code>path</code></td>
                <td>Validates PDF header and text layer.</td>
              </tr>
              <tr>
                <td><code>get_manifest</code></td>
                <td><code>folder_path</code></td>
                <td>Parses machine-readable AGENTS.md.</td>
              </tr>
            </tbody>
          </table>

          <div className="pixel-section-header">
            <span>IDE CONFIGURATION GENERATOR</span>
          </div>

          <div style={{ display: 'flex', gap: '6px', marginBottom: '8px', flexWrap: 'wrap' }}>
            {[
              { id: 'cursor', label: 'CURSOR' },
              { id: 'claude', label: 'CLAUDE_DESKTOP' },
              { id: 'antigravity', label: 'ANTIGRAVITY' },
              { id: 'vscode', label: 'VS_CODE' },
              { id: 'zed', label: 'ZED_WINDSURF' }
            ].map((t) => (
              <button
                key={t.id}
                className={`pixel-btn ${ideTab === t.id ? 'active' : ''}`}
                onClick={() => setIdeTab(t.id)}
                style={{ fontSize: '11px', padding: '3px 8px' }}
              >
                [{t.label}]
              </button>
            ))}
          </div>

          <CodeBlock code={configs[ideTab]} language="json" title={`${ideTab.toUpperCase()} CONFIG`} />
        </div>
      </div>
    </div>
  );
}
