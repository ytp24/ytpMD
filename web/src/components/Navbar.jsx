import React from 'react';

export function Navbar({ currentPage, setPage }) {
  const pages = [
    { id: 'home', num: '01', label: 'HOME' },
    { id: 'docs', num: '02', label: 'DOCS' },
    { id: 'downloads', num: '03', label: 'DOWNLOAD' },
    { id: 'mcp', num: '04', label: 'MCP_SPEC' },
    { id: 'legal', num: '05', label: 'LEGAL' }
  ];

  return (
    <nav className="pixel-nav">
      {pages.map((p) => (
        <button
          key={p.id}
          className={`pixel-nav-item ${currentPage === p.id ? 'active' : ''}`}
          onClick={() => setPage(p.id)}
        >
          [{p.num}:{p.label}]
        </button>
      ))}
    </nav>
  );
}

export function Footer() {
  return (
    <footer className="pixel-footer">
      <div>
        <span style={{ color: '#2dd4bf' }}>&gt; STATUS:</span> ONLINE &bull; 
        <span style={{ color: '#2dd4bf' }}> TELEMETRY:</span> 0% &bull; 
        <span style={{ color: '#2dd4bf' }}> ENGINE:</span> GO 1.22+
      </div>
      <div>
        <a
          href="https://github.com/ytp24/ytpMD"
          target="_blank"
          rel="noreferrer"
          style={{ color: '#14b8a6', textDecoration: 'none' }}
        >
          [ GITHUB_REPO &nearr; ]
        </a>
      </div>
    </footer>
  );
}
