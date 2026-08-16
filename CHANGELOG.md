# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [2.0.0] - 2026-08-16

### 🚀 Major Rewrite & Features
- **Golang Rewrite**: Complete rewrite of the Nyaa & Sukebei API service from Python (FastAPI/BeautifulSoup4) to Golang (`net/http`, `goquery`).
- **High Performance & Low Memory Footprint**: Drastically reduced RAM usage (~15MB vs ~100MB+) and improved request throughput and response latency.
- **Static Binary & Multi-Stage Docker**: Updated `Dockerfile` to produce a minimal, lightweight Alpine image with multi-stage Go build.

### ⚡ Improvements & Quality Automation
- **Go Linting**: Added `golangci-lint` GitHub Actions workflow (`.github/workflows/golangci-lint.yml`).
- **SonarQube Integration**: Added SonarQube code analysis workflow (`.github/workflows/sonar.yml` and `sonar-project.properties`).
- **Tag-Based Docker Releases**: Enhanced Docker GitHub Actions workflow (`.github/workflows/docker-build-push.yml`) to support semver tag-based builds (`v*.*.*`) pushing multi-tagged images to GitHub Container Registry (GHCR).

### 🐛 Fixes
- Removed legacy Python dependencies and virtualenv overhead.
