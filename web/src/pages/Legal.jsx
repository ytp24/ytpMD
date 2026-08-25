import React from 'react';

export default function Legal() {
  return (
    <div>
      <div className="pixel-box">
        <div className="pixel-box-header">
          <span>/// 05_LEGAL_COMPLIANCE</span>
          <span>&gt; POLICY</span>
        </div>

        <div className="pixel-box-body">
          <div className="pixel-section-header">
            <span>[1.0] APACHE 2.0 OPEN SOURCE LICENSE</span>
          </div>
          <p style={{ color: '#cbd5e1', marginBottom: '14px' }}>
            ytpMD is open-source software licensed under the <strong>Apache License, Version 2.0</strong>.
            Copyright (c) 2026 ytp24 (<code>ykinwork24@gmail.com</code>).
          </p>

          <div className="pixel-section-header" style={{ marginTop: '24px' }}>
            <span>[2.0] ZERO-TELEMETRY GUARANTEE</span>
          </div>
          <ul style={{ paddingLeft: '18px', color: '#cbd5e1', lineHeight: '1.8' }}>
            <li><strong>100% Local-First:</strong> All parsing and transformation run on your local machine.</li>
            <li><strong>Zero Tracking:</strong> No telemetry, analytics, IP logging, or tracking tokens.</li>
            <li><strong>Zero Cloud Dependency:</strong> Fully operational in air-gapped environments.</li>
          </ul>

          <div className="pixel-section-header" style={{ marginTop: '24px' }}>
            <span>[3.0] SECURITY DISCLOSURE SLA</span>
          </div>
          <table className="pixel-table">
            <thead>
              <tr>
                <th>Phase</th>
                <th>SLA Target</th>
              </tr>
            </thead>
            <tbody>
              <tr>
                <td>Acknowledgment</td>
                <td>Within 24 hours</td>
              </tr>
              <tr>
                <td>Triage & Severity Assessment</td>
                <td>Within 48 hours</td>
              </tr>
              <tr>
                <td>Patch Release</td>
                <td>Within 7 business days</td>
              </tr>
            </tbody>
          </table>
          <p style={{ fontSize: '11.5px', color: '#94a3b8', marginTop: '8px' }}>
            Report vulnerabilities to <code>ykinwork24@gmail.com</code> with subject <code>[SECURITY] Vulnerability Report</code>.
          </p>

          <div className="pixel-section-header" style={{ marginTop: '24px' }}>
            <span>[4.0] SEMVER LIFECYCLE (v1.0.0 &rarr; v3.2.0)</span>
          </div>
          <table className="pixel-table">
            <thead>
              <tr>
                <th>Version</th>
                <th>Status</th>
                <th>Key Capabilities</th>
              </tr>
            </thead>
            <tbody>
              <tr>
                <td><code>v3.2.0</code></td>
                <td>Active Stable</td>
                <td>MCP Server stdio, Supabase binary storage, Snap/Brew/AUR packages.</td>
              </tr>
              <tr>
                <td><code>v3.1.0</code></td>
                <td>Maintenance</td>
                <td>Concurrent batch worker pool, ANSI progress bar, Clean architecture.</td>
              </tr>
              <tr>
                <td><code>v2.0.0</code></td>
                <td>Maintenance</td>
                <td>Interactive prompt wizard, native OS GUI dialog fallbacks.</td>
              </tr>
              <tr>
                <td><code>v1.0.0</code></td>
                <td>Historical</td>
                <td>Initial Go release, Appendix/Index cutoff filter, paragraph reflower.</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}
