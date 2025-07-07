# CI/CD Workflows

This repository has two main workflows:

## 🚀 CI/CD Pipeline (`ci.yml`)

**Triggers:**

- Push to main or develop branches
- Pull requests to main branch
- Manual trigger from Actions tab

**Purpose:** Continuous integration, testing, and container management

```mermaid
flowchart TD
    A[test] --> B[build]
    A --> C[docker-build]
    C --> D[cleanup-old-images]
    F[security-scan]

    A -.-> |"Always runs first"| A
    C -.-> |"Skipped on PRs"| C
    D -.-> |"Skipped on PRs"| D
    F -.-> |"Runs independently, skipped on PRs"| F

    style A fill:#e1f5fe
    style B fill:#f3e5f5
    style C fill:#fff3e0
    style D fill:#fff3e0
    style F fill:#fce4ec
```

## 🏷️ Create Release (`release.yml`)

**Triggers:**

- Manual trigger only (workflow_dispatch)

**Purpose:** Semantic versioning and release management

```mermaid
flowchart TD
    G[get-latest-tag] --> H[calculate-next-version]
    H --> I[generate-release-notes]
    I --> J[run-tests]
    J --> K[build-cross-platform-binaries]
    K --> L[create-github-release]

    G -.-> |"Fetch version history"| G
    H -.-> |"Semantic versioning"| H
    I -.-> |"Auto-generate with PRs & contributors"| I
    J -.-> |"Ensure quality before release"| J
    K -.-> |"Linux, macOS, Windows"| K
    L -.-> |"Upload binaries & checksums"| L

    style G fill:#e8f5e8
    style H fill:#e8f5e8
    style I fill:#e8f5e8
    style J fill:#e8f5e8
    style K fill:#e8f5e8
    style L fill:#e8f5e8
```

## Usage

1. **Development**: Push code → CI runs tests, builds, and manages containers
2. **Release**: Go to Actions → Run "🏷️ Create Release" → Choose version type (patch/minor/major)
