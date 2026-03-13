# Contributing to gpkg

Thanks for contributing! This document explains the basic workflow and expectations.

Getting started
1. Fork the repository and create a branch for your work:
   git checkout -b feat/my-feature
2. Make small, focused commits with clear messages.
3. Run tests and linters locally before opening a PR.

Issue workflow
- Pick an open issue (start with an MVP task or ask for assignment).
- If you need the issue split, request it in the issue comments or ask in the Space.
- Link your PR to the issue with "Closes #<issue>".

Coding standards
- Language: Go (use gofmt/golangci-lint)
- Include unit tests for new logic where feasible.
- Keep the CLI ergonomics predictable and script-friendly. Provide --json outputs for machine consumption.

Documentation
- Update README.md and examples when adding features.
- Add example manifests for new manifest-supported cases.

Communication
- Use concise PR descriptions and explain any backward-incompatible changes.
- For design decisions, create an issue or RFC and discuss before large work items.

If you want help picking a starter issue, say "suggest starter tasks" and I will list a few good first PRs.