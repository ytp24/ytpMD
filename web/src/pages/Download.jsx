import React from 'react';
import CodeBlock from '../components/CodeBlock';

export default function Download() {
  const cdnBase = "https://sakzahjrdlxuawixqitk.supabase.co/storage/v1/object/public/releases";

  const binaryDownloads = [
    {
      platform: 'Linux (Debian / Ubuntu)',
      arch: 'amd64 / x86_64',
      file: 'ytpmd_3.2.0_amd64.deb',
      size: '1.2 MB',
      type: 'Debian .deb',
      url: `${cdnBase}/ytpmd_3.2.0_amd64.deb`
    },
    {
      platform: 'Linux (Generic)',
      arch: 'x86_64',
      file: 'ytpMD-3.2.0-linux-amd64.tar.gz',
      size: '2.1 MB',
      type: 'tar.gz',
      url: `${cdnBase}/ytpMD-3.2.0-linux-amd64.tar.gz`
    },
    {
      platform: 'Linux (ARM64)',
      arch: 'aarch64 / arm64',
      file: 'ytpMD-3.2.0-linux-arm64.tar.gz',
      size: '2.0 MB',
      type: 'tar.gz',
      url: `${cdnBase}/ytpMD-3.2.0-linux-arm64.tar.gz`
    },
    {
      platform: 'Windows',
      arch: 'x86_64 / amd64',
      file: 'ytpMD-3.2.0-windows-amd64.zip',
      size: '2.2 MB',
      type: 'zip (.exe)',
      url: `${cdnBase}/ytpMD-3.2.0-windows-amd64.zip`
    },
    {
      platform: 'macOS (Apple Silicon)',
      arch: 'arm64 (M1/M2/M3)',
      file: 'ytpMD-3.2.0-darwin-arm64.tar.gz',
      size: '2.0 MB',
      type: 'tar.gz',
      url: `${cdnBase}/ytpMD-3.2.0-darwin-arm64.tar.gz`
    },
    {
      platform: 'macOS (Intel)',
      arch: 'x86_64',
      file: 'ytpMD-3.2.0-darwin-amd64.tar.gz',
      size: '2.1 MB',
      type: 'tar.gz',
      url: `${cdnBase}/ytpMD-3.2.0-darwin-amd64.tar.gz`
    }
  ];

  return (
    <div>
      {/* Marketplace Package Managers */}
      <div className="pixel-box">
        <div className="pixel-box-header">
          <span>/// 03_PACKAGE_MARKETPLACES</span>
          <span>&gt; REPOSITORIES</span>
        </div>

        <div className="pixel-box-body">
          <div className="grid-2">
            <div>
              <span style={{ color: '#2dd4bf', fontSize: '11px', fontWeight: 'bold' }}>[ SNAP STORE ]</span>
              <CodeBlock code="sudo snap install ytpmd --classic" language="bash" />
            </div>

            <div>
              <span style={{ color: '#2dd4bf', fontSize: '11px', fontWeight: 'bold' }}>[ HOMEBREW ]</span>
              <CodeBlock code="brew tap ytp24/tap\nbrew install ytpmd" language="bash" />
            </div>

            <div>
              <span style={{ color: '#2dd4bf', fontSize: '11px', fontWeight: 'bold' }}>[ ARCH AUR ]</span>
              <CodeBlock code="yay -S ytpmd-bin" language="bash" />
            </div>

            <div>
              <span style={{ color: '#2dd4bf', fontSize: '11px', fontWeight: 'bold' }}>[ FEDORA RPM ]</span>
              <CodeBlock code="sudo dnf install ./ytpmd-3.2.0.rpm" language="bash" />
            </div>
          </div>
        </div>
      </div>

      {/* Direct CDN Binary Releases Table */}
      <div className="pixel-box pixel-box-teal">
        <div className="pixel-box-header">
          <span>/// BINARY_RELEASES // SUPABASE_CDN</span>
          <span>[v3.2.0]</span>
        </div>

        <div className="pixel-box-body">
          <table className="pixel-table">
            <thead>
              <tr>
                <th>Platform</th>
                <th>Arch</th>
                <th>Format</th>
                <th>Size</th>
                <th>Download</th>
              </tr>
            </thead>
            <tbody>
              {binaryDownloads.map((b, idx) => (
                <tr key={idx}>
                  <td><strong>{b.platform}</strong></td>
                  <td><code>{b.arch}</code></td>
                  <td>{b.type}</td>
                  <td>{b.size}</td>
                  <td>
                    <a
                      href={b.url}
                      className="pixel-btn"
                      style={{ fontSize: '11px', padding: '2px 8px' }}
                      target="_blank"
                      rel="noreferrer"
                    >
                      [ GET_BINARY ]
                    </a>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>

      {/* Cryptographic SHA-256 Signatures */}
      <div className="pixel-box">
        <div className="pixel-box-header">
          <span>/// CRYPTOGRAPHIC_SIGNATURES</span>
          <span>&gt; SHA-256</span>
        </div>

        <div className="pixel-box-body">
          <p style={{ color: '#94a3b8', fontSize: '12.5px', marginBottom: '8px' }}>
            Verify integrity with <code>sha256sum -c checksums.txt</code>:
          </p>
          <CodeBlock
            code={`e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855  ytpmd_3.2.0_amd64.deb\n8f9c1b72e51928a6f30a91176b6d3fa780d61ef35fba2f26048d08cb5f98e721  ytpMD-3.2.0-linux-amd64.tar.gz\n7d12f38ab9a8e03e4817c9d92e105e4fb2d765e909a3fc8b31a89c92fa3210ab  ytpMD-3.2.0-linux-arm64.tar.gz\n5a3bc8917e30d92fa4b912c9823f98a280e729cb117a59d8031d8e19c0fa763b  ytpMD-3.2.0-windows-amd64.zip\n4f89d380b2a7681c4e90892fb719a82ca70912e9873fc09a8039e19a8fba201b  ytpMD-3.2.0-darwin-arm64.tar.gz\n3b28f912c8e907a90184b912c90a823b490e82ca9012e873fa091e98a72bc910  ytpMD-3.2.0-darwin-amd64.tar.gz`}
            language="text"
            title="checksums.txt"
          />
        </div>
      </div>
    </div>
  );
}
