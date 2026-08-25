Name:           ytpmd
Version:        3.2.0
Release:        1%{?dist}
Summary:        High-Performance Concurrent PDF to Chapter Markdown Engine & MCP Server
License:        Apache-2.0
URL:            https://github.com/ytp24/ytpMD
Requires:       poppler-utils

%description
ytpMD [pdf2md] transforms PDF documents into clean, structured GitHub-Flavored
Markdown notes grouped by chapter with zero noise, automatically excluding
appendix and non-usable sections, with full AI agent manifest and MCP support.

%install
mkdir -p %{buildroot}%{_bindir}
install -m 755 ytpmd %{buildroot}%{_bindir}/ytpmd
install -m 755 ytpmd-mcp %{buildroot}%{_bindir}/ytpmd-mcp
ln -sf ytpmd %{buildroot}%{_bindir}/ytpMD
ln -sf ytpmd %{buildroot}%{_bindir}/ytp24
ln -sf ytpmd %{buildroot}%{_bindir}/pdf2md

%files
%{_bindir}/ytpmd
%{_bindir}/ytpmd-mcp
%{_bindir}/ytpMD
%{_bindir}/ytp24
%{_bindir}/pdf2md

%changelog
* Wed Aug 26 2026 ytp24 <ykinwork24@gmail.com> - 3.2.0-1
- Model Context Protocol (MCP) server integration
- Concurrent batch worker pool
- Cross-platform packaging support
