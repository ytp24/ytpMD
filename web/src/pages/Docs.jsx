import React from 'react';
import CodeBlock from '../components/CodeBlock';

export default function Docs() {
  return (
    <div>
      <div className="pixel-box">
        <div className="pixel-box-header">
          <span>/// 02_DOCUMENTATION // MANUAL</span>
          <span>&gt; CLI_SPEC</span>
        </div>

        <div className="pixel-box-body">
          <div className="pixel-section-header">
            <span>[1.0] PREREQUISITES</span>
          </div>
          <p style={{ color: '#cbd5e1', marginBottom: '8px' }}>
            ytpMD requires the system package <code>poppler-utils</code> for text stream extraction:
          </p>
          <CodeBlock
            code={`# Ubuntu / Debian:\nsudo apt install poppler-utils\n\n# Fedora / RHEL:\nsudo dnf install poppler-utils\n\n# Arch Linux:\nsudo pacman -S poppler\n\n# macOS (Homebrew):\nbrew install poppler`}
            language="bash"
            title="INSTALL POPPLER-UTILS"
          />

          <div className="pixel-section-header" style={{ marginTop: '24px' }}>
            <span>[2.0] CLI COMMANDS</span>
          </div>

          <table className="pixel-table">
            <thead>
              <tr>
                <th>Command</th>
                <th>Arguments</th>
                <th>Action</th>
              </tr>
            </thead>
            <tbody>
              <tr>
                <td><code>ytpmd</code></td>
                <td><em>(none)</em></td>
                <td>Interactive wizard with system file picker dialog.</td>
              </tr>
              <tr>
                <td><code>ytpmd convert</code></td>
                <td><code>&lt;file.pdf&gt;</code></td>
                <td>Extract PDF into chapter folder or single Markdown file.</td>
              </tr>
              <tr>
                <td><code>ytpmd batch</code></td>
                <td><code>&lt;dir&gt;</code></td>
                <td>Concurrent batch convert with Goroutine worker pool.</td>
              </tr>
              <tr>
                <td><code>ytpmd mcp</code></td>
                <td><em>(none)</em></td>
                <td>Model Context Protocol JSON-RPC 2.0 stdio server for IDEs.</td>
              </tr>
              <tr>
                <td><code>ytpmd version</code></td>
                <td><em>(none)</em></td>
                <td>Display version string (<code>v3.2.0</code>) and Go runtime.</td>
              </tr>
            </tbody>
          </table>

          <div className="pixel-section-header" style={{ marginTop: '24px' }}>
            <span>[3.0] CLI FLAGS</span>
          </div>

          <table className="pixel-table">
            <thead>
              <tr>
                <th>Flag</th>
                <th>Type</th>
                <th>Default</th>
                <th>Purpose</th>
              </tr>
            </thead>
            <tbody>
              <tr>
                <td><code>-o, -output</code></td>
                <td>string</td>
                <td><code>~/Documents/ytpMD</code></td>
                <td>Root destination folder for extracted notes.</td>
              </tr>
              <tr>
                <td><code>-concurrency</code></td>
                <td>int</td>
                <td>4</td>
                <td>Parallel worker Goroutines in batch mode.</td>
              </tr>
              <tr>
                <td><code>-f, -force</code></td>
                <td>bool</td>
                <td>false</td>
                <td>Overwrite destination directory without prompting.</td>
              </tr>
              <tr>
                <td><code>-skip-front &lt;N&gt;</code></td>
                <td>int</td>
                <td>0</td>
                <td>Skip first N pages (covers, dedication, TOC).</td>
              </tr>
              <tr>
                <td><code>-keep-appendix</code></td>
                <td>bool</td>
                <td>false</td>
                <td>Do NOT stop extraction at Appendix or Index.</td>
              </tr>
              <tr>
                <td><code>-single-file</code></td>
                <td>bool</td>
                <td>false</td>
                <td>Output single monolithic .md instead of chapter folder.</td>
              </tr>
            </tbody>
          </table>

          <div className="pixel-section-header" style={{ marginTop: '24px' }}>
            <span>[4.0] AGENT MANIFEST & FRONTMATTER</span>
          </div>

          <CodeBlock
            code={`---
title: "CHAPTER 1: KUBERNETES CLUSTER ARCHITECTURE"
chapter: 1
total_chapters: 8
source_document: "DevOps_Handbook.pdf"
start_page: 14
word_count: 2450
estimated_tokens: 3200
agent_instructions: "Cite section headers and use code snippets directly when referencing."
---`}
            language="yaml"
            title="CHAPTER YAML HEADER"
          />
        </div>
      </div>
    </div>
  );
}
