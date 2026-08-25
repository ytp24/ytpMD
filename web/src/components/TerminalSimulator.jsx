import React, { useState } from 'react';

const SAMPLES = {
  wizard: `
  ██╗   ██╗████████╗██████╗ ███╗   ███╗██████╗ 
  ╚██╗ ██╔╝╚══██╔══╝██╔══██╗████╗ ████║██╔══██╗
   ╚████╔╝    ██║   ██████╔╝██╔████╔██║██║  ██║
    ╚██╔╝     ██║   ██╔═══╝ ██║╚██╔╝██║██║  ██║
     ██║      ██║   ██║     ██║ ╚═╝ ██║██████╔╝
     ╚═╝      ╚═╝   ╚═╝     ╚═╝     ╚═╝╚═════╝ 
               [ pdf2md ]  v3.2.0
   -- High-Performance Concurrent PDF to Markdown Engine --
   Transforms PDFs into Chapter-Aware Markdown Notes with Zero Noise

[?] Select Mode:
    1. Single PDF File (default)
    2. Batch Directory of PDFs (Concurrent Goroutines)
Choice [1]: 1

[?] Enter PDF path [or press Enter for GUI dialog]: /docs/Cloud_Architecture.pdf
[+] Selected file: /docs/Cloud_Architecture.pdf
[?] Enter destination root [~/Documents/ytpMD]: 
[?] Apply production defaults (TOC chapters -> Appendix cutoff)? (Y/n): Y

[*] Extracting and transforming: Cloud_Architecture.pdf...
[+] Complete (Agent-Ready):
   ├── Folder:   ~/Documents/ytpMD/Cloud_Architecture/
   ├── Index:    ~/Documents/ytpMD/Cloud_Architecture/README.md
   ├── Manifest: ~/Documents/ytpMD/Cloud_Architecture/AGENTS.md
   └── Chapters: 4 extracted (~13,670 tokens total)
`,
  convert: `
ytpmd@local:~# ytpmd convert Kubernetes_DeepDive.pdf -force

[*] Processing: Kubernetes_DeepDive.pdf...
[+] Extraction Complete -> ~/Documents/ytpMD/Kubernetes_DeepDive/
   ├── [01] Container Runtimes         -> 01_container_runtimes.md (~3,100 tokens)
   ├── [02] Control Plane Architecture -> 02_control_plane.md (~4,200 tokens)
   ├── [03] Custom Resources (CRDs)    -> 03_crds.md (~2,900 tokens)
   ├── [04] Ingress & Networking       -> 04_ingress.md (~3,400 tokens)
   ├── [05] Service Mesh Topology      -> 05_service_mesh.md (~4,600 tokens)
   └── [06] Observability Pipelines    -> 06_observability.md (~3,800 tokens)
   └── Metrics: 312 pages total | 278 converted | 34 noise pages truncated.
`,
  batch: `
ytpmd@local:~# ytpmd batch ~/Books/DevOps/ -name DevOpsLib -concurrency 6 -force

[*] Launching concurrent batch engine (6 Goroutines) for 4 PDF(s)...

Converting Batch [████████████████████████████] 100% (4/4) | Complete

[+] Batch Finished in 1.482s -> ~/Documents/ytpMD/DevOpsLib/
   ├── Cloud_Architecture.pdf     (4 chapters, 148 pages)
   ├── Kubernetes_DeepDive.pdf    (6 chapters, 312 pages)
   ├── Linux_Kernel_Internals.pdf (8 chapters, 420 pages)
   └── Terraform_Mastery.pdf      (5 chapters, 210 pages)
`,
  mcp: `
ytpmd@local:~# ytpmd mcp

[ytpMD MCP] Stdio JSON-RPC 2.0 Server Active (v3.2.0)
--> JSON-RPC 'initialize' received from IDE
<-- Server Capabilities: tools={convert_pdf, batch_convert, inspect_pdf, get_manifest}
--> 'tools/call' convert_pdf(path="DevOps_Guide.pdf")
[*] Extracted 8 chapters into ~/Documents/ytpMD/DevOps_Guide/
<-- Returned chapter index, token map, and AGENTS.md manifest to AI agent
`
};

export default function TerminalSimulator() {
  const [activeTab, setActiveTab] = useState('wizard');

  return (
    <div className="pixel-box" style={{ margin: '20px 0' }}>
      <div className="pixel-box-header">
        <span>&gt; CONSOLE_OUTPUT // ytpMD v3.2.0</span>
        <span style={{ color: '#4ade80' }}>● ONLINE</span>
      </div>

      <div style={{ padding: '8px 12px', backgroundColor: '#090e17', borderBottom: '1px solid #1e293b', display: 'flex', gap: '6px', flexWrap: 'wrap' }}>
        {[
          { id: 'wizard', label: '1: WIZARD' },
          { id: 'convert', label: '2: CONVERT' },
          { id: 'batch', label: '3: BATCH' },
          { id: 'mcp', label: '4: MCP_SERVER' }
        ].map((t) => (
          <button
            key={t.id}
            className={`pixel-btn ${activeTab === t.id ? 'active' : ''}`}
            onClick={() => setActiveTab(t.id)}
            style={{ fontSize: '11px', padding: '3px 8px' }}
          >
            [{t.label}]
          </button>
        ))}
      </div>

      <div className="pixel-term">
        <pre style={{ margin: 0, whiteSpace: 'pre', overflowX: 'auto' }}>
          {SAMPLES[activeTab]}
        </pre>
      </div>
    </div>
  );
}
