# ThousandEyes Synthetics CI

Validate HTTP endpoints, browser page loads, and network connectivity after a change with ThousandEyes instant tests. Run the same suite from GitHub Actions, Jenkins, GitLab CI, a developer terminal, Argo CD, or Argo Rollouts.

The Go runner uses the official [ThousandEyes CLI](https://github.com/thousandeyes/thousandeyes-cli) for API access, evaluates per-agent thresholds, and produces JSON/JUnit reports. The GitHub Action adds pipeline outputs and a job summary.

## Choose your platform

Each integration runs `te-synthetics run-tests` with the same version 1 suite.

| Platform | Setup guide | Example |
| --- | --- | --- |
| GitHub Actions | [Action setup](docs/platforms.md#github-actions) | [Workflow](examples/github-actions/ci.yml) |
| Developer terminal | [Native command](docs/platforms.md#developer-terminal) | [Suite](examples/synthetics/thousandeyes-ci.json) |
| Jenkins | [Worker and credentials](docs/platforms.md#jenkins) | [Jenkinsfile](examples/jenkins/Jenkinsfile) |
| GitLab CI | [Runner and variables](docs/platforms.md#gitlab-ci) | [Pipeline](examples/gitlab/ci.yml) |
| Container | [Build and run](docs/platforms.md#container) | [Dockerfile](Dockerfile) |
| Argo CD | [PostSync validation](docs/platforms.md#argo-cd) | [Job](examples/argocd/post-sync-job.yaml) |
| Argo Rollouts | [Progressive delivery](docs/platforms.md#argo-rollouts) | [AnalysisTemplate](examples/synthetics/argo-rollouts/analysis-template.yaml) and [rollout step](examples/synthetics/argo-rollouts/rollout-step.yaml) |

## Configure your checks

Copy the [example suite](examples/synthetics/thousandeyes-ci.json), then replace its account group, agents, targets, and thresholds. Supply `TE_TOKEN` through your secret store or shell environment. Private targets require Enterprise Agents that can reach them. Each run creates instant tests and consumes applicable ThousandEyes units.

Use the [JSON Schema](schemas/thousandeyes-ci.schema.json) for editor completion. See [configuration and results](docs/configuration.md) for command options, thresholds, exit codes, and report semantics. Start with the [Page Load example](docs/configuration.md#page-load-example) to measure browser performance.

## Development and releases

- [DEVELOPMENT.md](DEVELOPMENT.md): prerequisites, source builds, coding conventions, and tests.
- [CONTRIBUTING.md](CONTRIBUTING.md): contribution policy and pull requests.
- [RELEASING.md](RELEASING.md): version pinning, dependencies, and publication.
- [Architecture](docs/architecture.md): components and the official CLI backend contract.

## Community and support

Use [GitHub issues](https://github.com/thousandeyes/thousandeyes-synthetics-ci/issues) for bugs, questions, and feature requests. Check existing discussions first, and remove credentials and private account or target information from examples and reports.

Follow the [code of conduct](CODE_OF_CONDUCT.md). Report vulnerabilities through the private channel in [SECURITY.md](SECURITY.md).

## License and attribution

Licensed under Apache-2.0. See [LICENSE](LICENSE), [NOTICE](NOTICE), and [third-party attribution](docs/licensing.md).
