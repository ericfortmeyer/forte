# Contributing to Forte

Thanks for your interest — contributions are welcome.

Quick start

1. Fork the repo and create a branch from main: git checkout -b feat/short-description
2. Run tests locally: make -s test; make -s docker-integration-test
3. Commit using a short prefix (feat/, fix/, docs/). Example: feat(deploy): support user namespaces
4. Open a pull request (see template).

PRs

- Include a short description, linked issue (if any), and testing steps.
- Add or update tests for code changes.
- Keep changes focused and small where possible.

Branching & versioning

- Branch from main.
- We use semantic versioning. Tag releases as MAJOR.MINOR.PATCH.

Review & merging

- The maintainer reviews and merges PRs. CI must pass before merge.
- Maintainer may request changes or suggest alternative approaches.

Reporting bugs & feature requests

- Open an issue using the template. Provide reproduction steps and environment.

How contributors are credited

- Contributors are listed in CONTRIBUTORS.md or via git history.
