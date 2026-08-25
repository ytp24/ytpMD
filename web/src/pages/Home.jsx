import React, { useState } from 'react';
import CodeBlock from '../components/CodeBlock';
import TerminalSimulator from '../components/TerminalSimulator';
import "./styles/logo.css"
import LogoSVG from "../components/svg/LogoSVG";
export default function Home({ setPage }) {
  const [installTab, setInstallTab] = useState('curl');

  const installCommands = {
    curl: 'curl -fsSL https://raw.githubusercontent.com/ytp24/ytpMD/main/scripts/install.sh | bash',
    wget: 'wget -qO- https://raw.githubusercontent.com/ytp24/ytpMD/main/scripts/install.sh | bash',
    powershell: 'irm https://raw.githubusercontent.com/ytp24/ytpMD/main/scripts/install.ps1 | iex',
    snap: 'sudo snap install ytpmd --classic',
    brew: 'brew install ytp24/tap/ytpmd',
    deb: 'sudo apt install ./dist/ytpmd_3.2.0_amd64.deb'
  };

  return (
    <div>
      {/* 1990s Minimalist Pixel Hero Header */}
      <div className="pixel-box pixel-box-teal" style={{ padding: '24px 20px', textAlign: 'center', marginBottom: '24px' }}>
<div className="ytpmd-logo">
  <LogoSVG />
</div>

        <div className="pixel-title" style={{ marginTop: '4px' }}>
          [ pdf2md ] v3.2.0
        </div>
        <p style={{ color: '#94a3b8', fontSize: '14px', maxWidth: '680px', margin: '10px auto 16px auto', lineHeight: '1.5' }}>
          High-performance, local-first engine converting technical PDF manuals into clean, chapter-segmented Markdown notes formatted with YAML frontmatter & AI agent manifests.
        </p>

        <div style={{ display: 'flex', justifyContent: 'center', gap: '6px', flexWrap: 'wrap', marginBottom: '18px' }}>
          <span className="pixel-badge">[ GO 1.22 ]</span>
          <span className="pixel-badge">[ MCP READY ]</span>
          <span className="pixel-badge pixel-badge-green">[ 100% OFFLINE ]</span>
          <span className="pixel-badge">[ ZERO TELEMETRY ]</span>
          <span className="pixel-badge pixel-badge-amber">[ APACHE 2.0 ]</span>
        </div>

        <div style={{ display: 'flex', justifyContent: 'center', gap: '10px', flexWrap: 'wrap' }}>
          <button className="pixel-btn pixel-btn-active" onClick={() => setPage('downloads')}>
            [ DOWNLOAD_BINARY ]
          </button>
          <button className="pixel-btn" onClick={() => setPage('mcp')}>
            [ CONFIGURE_MCP ]
          </button>
          <button className="pixel-btn" onClick={() => setPage('docs')}>
            [ READ_MANUAL ]
          </button>
        </div>
      </div>

      {/* One-Line Instant Installation */}
      <div className="pixel-box">
        <div className="pixel-box-header">
          <span>/// 01_QUICK_INSTALLATION</span>
          <span>&gt; ONE_LINER</span>
        </div>
        <div className="pixel-box-body">
          <div style={{ display: 'flex', gap: '6px', marginBottom: '8px', flexWrap: 'wrap' }}>
            {['curl', 'wget', 'powershell', 'snap', 'brew', 'deb'].map((tool) => (
              <button
                key={tool}
                className={`pixel-btn ${installTab === tool ? 'active' : ''}`}
                onClick={() => setInstallTab(tool)}
                style={{ fontSize: '11px', padding: '3px 8px', textTransform: 'uppercase' }}
              >
                [{tool}]
              </button>
            ))}
          </div>

          <CodeBlock code={installCommands[installTab]} language="bash" title={`INSTALL VIA ${installTab.toUpperCase()}`} />
        </div>
      </div>

      {/* Terminal Simulator Console */}
      <TerminalSimulator />

      {/* Minimalist 1990s Problem vs Solution */}
      <div className="grid-2">
        <div className="pixel-box">
          <div className="pixel-box-header" style={{ color: '#f87171' }}>
            <span>[X] MONOLITHIC PDF ISSUES</span>
          </div>
          <div className="pixel-box-body" style={{ fontSize: '12.5px', color: '#cbd5e1' }}>
            <ul style={{ paddingLeft: '18px', lineHeight: '1.8' }}>
              <li><strong>Token Limit Blowout:</strong> Cannot feed 400-page books into LLMs.</li>
              <li><strong>Context Noise:</strong> Polluted with page headers, image tags, and footers.</li>
              <li><strong>Broken Hyphenation:</strong> Words split mid-line (<code>archi- \n tecture</code>).</li>
              <li><strong>Deadweight Index:</strong> False-positive retrieval from back matter.</li>
            </ul>
          </div>
        </div>

        <div className="pixel-box pixel-box-teal">
          <div className="pixel-box-header" style={{ color: '#4ade80' }}>
            <span>[✓] ytpMD ARCHITECTURE</span>
          </div>
          <div className="pixel-box-body" style={{ fontSize: '12.5px', color: '#cbd5e1' }}>
            <ul style={{ paddingLeft: '18px', lineHeight: '1.8' }}>
              <li><strong>Chapter Segmentation:</strong> Split cleanly into individual notes.</li>
              <li><strong>AI Agent Manifest:</strong> Automatic <code>AGENTS.md</code> with token counts.</li>
              <li><strong>Appendix Cutoff:</strong> Auto-detection truncates non-usable back matter.</li>
              <li><strong>MCP Integration:</strong> Stdio JSON-RPC 2.0 server for Cursor & Claude.</li>
            </ul>
          </div>
        </div>
      </div>
    </div>
  );
}
