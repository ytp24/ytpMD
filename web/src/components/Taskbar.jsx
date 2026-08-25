import React from 'react';

export function Taskbar({ currentPage, setPage }) {
  const pages = [
    { id: 'home', num: '01', label: 'HOME' },
    { id: 'docs', num: '02', label: 'DOCS' },
    { id: 'downloads', num: '03', label: 'DOWNLOAD' },
    { id: 'mcp', num: '04', label: 'MCP_SPEC' },
    { id: 'legal', num: '05', label: 'LEGAL' }
  ];

  return (
    <div className="pixel-taskbar">
      {/* {pages.map((p) => (
        <button
          key={p.id}
          className={`pixel-taskbar-item ${currentPage === p.id ? 'active' : ''}`}
          onClick={() => setPage(p.id)}
        >
          [{p.num}:{p.label}]
        </button>
      ))} */}
    </div>
  );
}
