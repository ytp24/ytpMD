import React, { useState } from 'react';

export default function CodeBlock({ code, language = 'bash', title }) {
  const [copied, setCopied] = useState(false);

  const handleCopy = () => {
    navigator.clipboard.writeText(code);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <div style={{
      backgroundColor: '#04070d',
      border: '1px solid #1e293b',
      margin: '12px 0',
      position: 'relative'
    }}>
      <div style={{
        backgroundColor: '#0b111d',
        borderBottom: '1px solid #1e293b',
        padding: '6px 12px',
        fontSize: '11px',
        color: '#64748b',
        display: 'flex',
        justifyContent: 'space-between',
        alignItems: 'center'
      }}>
        <span>&gt;_ {title || language.toUpperCase()}</span>
        <button
          onClick={handleCopy}
          className="pixel-btn"
          style={{
            padding: '2px 8px',
            fontSize: '10px',
            backgroundColor: copied ? '#14b8a6' : 'transparent',
            color: copied ? '#000000' : '#2dd4bf',
            borderColor: '#14b8a6'
          }}
        >
          {copied ? '[✓ COPIED]' : '[ COPY ]'}
        </button>
      </div>
      <div style={{ padding: '12px 14px', overflowX: 'auto' }}>
        <pre style={{
          margin: 0,
          fontFamily: "'Fira Code', 'Courier New', monospace",
          fontSize: '12.5px',
          color: '#5eead4',
          lineHeight: '1.45',
          whiteSpace: 'pre'
        }}>
          <code>{code}</code>
        </pre>
      </div>
    </div>
  );
}
