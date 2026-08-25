import React, { useState } from 'react';
import { Navbar } from './components/Navbar';
import { Taskbar } from './components/Taskbar';
import Home from './pages/Home';
import Docs from './pages/Docs';
import Download from './pages/Download';
import MCP from './pages/MCP';
import Legal from './pages/Legal';

export default function App() {
  const [currentPage, setPage] = useState('home');

  const renderPage = () => {
    switch (currentPage) {
      case 'home':
        return <Home setPage={setPage} />;
      case 'docs':
        return <Docs />;
      case 'downloads':
        return <Download />;
      case 'mcp':
        return <MCP />;
      case 'legal':
        return <Legal />;
      default:
        return <Home setPage={setPage} />;
    }
  };

  return (
    <div className="container">
      {/* Outer Main Windows 95 Application Window */}
      <div className="win-outset" style={{ boxShadow: '8px 8px 24px rgba(0,0,0,0.4)', borderRadius: '2px' }}>
        {/* Master Window Title Bar */}
        <div className="win-titlebar">
          <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
            <span style={{ fontSize: '18px', gap: '20px'}}>⚡</span>
            <span>ytpMD [pdf2md] v3.2.0 — High-Performance PDF to Markdown Engine</span>
          </div>
          {/* <div className="win-titlebar-controls">
            <button className="win-btn win-titlebar-btn">_</button>
            <button className="win-btn win-titlebar-btn">□</button>
            <button className="win-btn win-titlebar-btn">×</button>
          </div> */}
        </div>

        {/* Menu Bar & Tabs */}
        <Navbar currentPage={currentPage} setPage={setPage} />

        {/* Main Content Viewport */}
        <div style={{ backgroundColor: '#c0c0c0', minHeight: '600px' }}>
          {renderPage()}
        </div>

        {/* Window Status Bar */}
        <div style={{
          backgroundColor: '#c0c0c0',
          borderTop: '2px solid #808080',
          padding: '4px 10px',
          display: 'flex',
          justifyContent: 'space-between',
          fontSize: '12px',
          color: '#334155'
        }}>
          <div><strong>Status:</strong> Ready • Memory: Clean • Zero Telemetry</div>
          <div><strong>Go Engine:</strong> 1.22+ • Apache 2.0</div>
        </div>
      </div>

      {/* Retro Bottom Taskbar */}
      <Taskbar currentPage={currentPage} setPage={setPage} />
    </div>
  );
}
