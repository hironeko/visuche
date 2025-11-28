# visuche 🎯

**visuche** (visualization check) is a powerful GitHub repository analytics tool that provides insights on PR metrics, lead times, release cadence, code review patterns, and CI/CD performance.

## ✨ Features

- **📊 Pull Request Analytics**: Lead time (avg/median), review time, merge wait, approval→merge
- **🚀 Release Cadence**: Counts merges into `main/master` as releases (bots/Dependabot excluded)
- **💬 Code Review Insights**: Review comments + approvals are counted for coverage/quality
- **🚀 CI/CD Performance**: GitHub Actions workflow analysis
- **🌐 Bilingual Output**: `--lang en|jp` / `--jp` for Japanese output
- **⚡ Fast & Efficient**: Parallel fetching, chunked date ranges, smart sampling

## 🚀 Quick Start

```bash
# Interactive mode (recommended)
visuche

# Analyze specific repository
visuche --repo owner/repo --since 2024-01-01 --until 2024-12-31

# Japanese output
visuche --repo owner/repo --since 2024-01-01 --until 2024-12-31 --jp

# Analyze GitHub Actions
visuche actions --repo owner/repo --since 2024-01-01
```

## 📦 Installation

### One-liner Install (All Platforms)

```bash
curl -fsSL https://raw.githubusercontent.com/hironeko/visuche/main/install.sh | bash
```

This will automatically:
- Install Go (if missing)
- Install GitHub CLI (if missing) 
- Build and install `visuche` to `~/bin/`
- Configure your PATH

### Alternative Methods

#### Homebrew (macOS/Linux)

```bash
brew tap hironeko/visuche
brew install visuche
```

#### GitHub Releases

Download the latest binary for your platform from [GitHub Releases](https://github.com/hironeko/visuche/releases).

#### Build from Source

```bash
git clone https://github.com/hironeko/visuche.git
cd visuche
./install.sh
```

## 📋 Prerequisites

The install script will automatically handle dependencies, but you can install manually:

- Go 1.19+ - [Download here](https://golang.org/dl/)
- [GitHub CLI (`gh`)](https://cli.github.com/) - For GitHub API access
- Authenticated GitHub CLI session: `gh auth login`

## 📊 Sample Output

```
📊 Pull Request Statistics
===================================================

🔢 Basic Metrics:
| Total PRs                     | 134 |
| Merged PRs                    | 132 |
| Releases (main/master merges) | 120 |
| Merge Rate                    | 98.5% |

⏱️ Timing Metrics:
| METRIC                 | AVERAGE | MEDIAN |
| Lead Time              | 10h28m  | 24m    |
| Review Time            | 2h47m   | -      |
| Merge Wait Time        | 13h41m  | 5h     |
| Approval→Merge Time    | 6h12m   | 2h     |

💬 Code Review Analysis:
| Review Comments per PR | 0.2 | 0.0 | 8 |
| Review Coverage        | 14 PRs (10.4%) |
```

## 🎯 Use Cases

- **Development Teams**: Track team velocity and code review effectiveness
- **Engineering Managers**: Monitor DORA metrics and development health
- **DevOps Engineers**: Analyze CI/CD performance and failure patterns
- **Open Source Maintainers**: Understand contributor patterns and project health

## 🤔 Why visuche?

Unlike other repository analytics tools, visuche:

- **Zero Configuration**: Works out of the box with GitHub CLI
- **Fast Analysis**: Smart sampling for large repositories
- **Practical Metrics**: Focus on actionable insights
- **Local First**: No data sent to external services
- **Developer Friendly**: Built by developers, for developers

## 📖 Commands

### PR Analysis (Default)

```bash
visuche [flags]
```

**Flags:**
- `--repo string`: Repository in 'owner/repo' format
- `--since string`: Analyze PRs since date (YYYY-MM-DD)
- `--until string`: Analyze PRs until date (YYYY-MM-DD)
- `--author string`: Filter by author username
- `--label string`: Filter by label name
- `--lang string`: `en` (default) or `jp` for Japanese output
- `--jp`: Shortcut for `--lang jp`

### GitHub Actions Analysis

```bash
visuche actions [flags]
```

Analyzes CI/CD performance, workflow success rates, and failure patterns.

## 🔧 Advanced Usage

### Large Repositories

visuche automatically optimizes for large repositories using:
- Smart sampling (recent + distributed historical PRs)
- Parallel processing for date ranges
- GraphQL complexity management

### Custom Time Ranges

```bash
# Last quarter analysis
visuche --since 2024-10-01 --until 2024-12-31

# Specific team member
visuche --author username --since 2024-01-01

# Feature branch analysis
visuche --label "feature" --since 2024-06-01
```

## 🤝 Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## 📄 License

MIT License - see [LICENSE](LICENSE) for details.

## 🙏 Acknowledgments

- Inspired by [peco](https://github.com/peco/peco) for terminal UX
- Built on [GitHub CLI](https://cli.github.com/) for robust GitHub API access
- Thanks to the Go community for excellent CLI libraries

---

**Made with ❤️ for developers who love data-driven insights**
