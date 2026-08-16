# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [2.0.0] - 2026-08-16

### 🚀 Major Rewrite & Features
- **Golang Migration**: Complete rewrite of the Nyaa & Sukebei API service from Python (FastAPI/BeautifulSoup4) to Golang (`net/http`, `goquery`).
- **Interactive Swagger UI**: Integrated Swagger/OpenAPI documentation automatically generated via `swag` and hosted interactively at `/docs`.
- **Repository Renaming**: Renamed module, Docker images, and repository references to **`NyaaSi-API-Go`**.
- **High Performance & Low Footprint**: Reduced memory footprint to ~10MB RAM with sub-millisecond routing overhead.
- **Multi-Stage Docker**: Multi-architecture Alpine runtime (`amd64`/`arm64`) with non-root security compliance (`USER appuser`).

### 🏗️ Architecture & Documentation
- **System Architecture Diagrams**: Added Mermaid component topology, sequence flow diagrams, and module responsibility breakdowns to `README.md`.
- **API Documentation**: Detailed query parameters, response structures, and taxonomy mapping for both `nyaa.si` and `sukebei.nyaa.si`.

### ⚡ Code Quality & CI/CD
- **SonarCloud Quality Gate**: Resolved 100% of SonarCloud issues (Cognitive Complexity, Duplicate Code, String Constants, Docker Non-Root User).
- **High Test Coverage**: Comprehensive unit test suite (`main_test.go`, `pkg/scraper/api_test.go`, `pkg/scraper/parsers_test.go`) achieving **>85% overall statement coverage** (91.4% in `pkg/scraper`).
- **Automated Workflows**: GitHub Actions workflows for Go linting (`golangci-lint`), SonarQube analysis, and SemVer tag releases (`v*`).

### 🐛 Fixes
- **CORS Support**: Added global CORS middleware supporting `OPTIONS` preflight requests and `Access-Control-Allow-Origin: *`.
- **Endpoint Parity**: Corrected `handleGetSukebeiID` handler call to target `sukebeiAPI` scraper.
