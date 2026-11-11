# Contributing to ThousandEyes Synthetics CI

Contributions can include bug reports, documentation, tests, platform examples, and code. Follow the [code of conduct](CODE_OF_CONDUCT.md) in issues, pull requests, and other project discussions.

## Report issues and suggest features

Search the repository's issues and pull requests before opening a new issue. For a bug, include the runner version or source commit, operating system, integration platform, reproduction steps, and expected and actual results. Include a minimal suite and sanitized diagnostics when useful.

For a feature, describe the problem, intended users, expected behavior, and alternatives. Discuss changes to the suite format, backend contract, or supported test types before investing in a large implementation.

Use synthetic values such as `https://example.com` in shared examples. Remove tokens, credentials, private targets, internal domains, account identifiers, and agent details from reports and logs. Report suspected vulnerabilities using [SECURITY.md](SECURITY.md), rather than a public issue or pull request.

## Submit a pull request

1. Fork the repository and create a branch from the current `main`.
2. Keep the change focused and explain the problem it solves. Link a related issue when one exists.
3. Use a concise, imperative commit subject, such as `Document agent selection`.
4. Add or update meaningful tests for behavior changes and documentation for user-facing changes. Follow [DEVELOPMENT.md](DEVELOPMENT.md) for setup, conventions, and validation.
5. Complete the pull request template with the resulting behavior, checks performed, and any limitations or compatibility effects.
6. Open the pull request against `main` and respond to review feedback.

Do not commit real credentials or reports containing private account data. Preserve copyright notices and third-party attribution. Identify external code's source, license, and modifications in the pull request; see [licensing](docs/licensing.md).

By submitting a contribution for inclusion, you agree to the contribution terms in section 5 of the [Apache-2.0 license](LICENSE).

Release publication is a maintainer responsibility. See [RELEASING.md](RELEASING.md) for the release process.
