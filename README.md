# ThousandEyes Synthetics CI

Validate application and network behavior after a change with ThousandEyes HTTP-server and agent-to-server instant tests. Run the same suite from GitHub Actions, Jenkins, GitLab CI, a developer terminal, Argo CD, or Argo Rollouts.

The Go runner calls the official [ThousandEyes CLI](https://github.com/thousandeyes/thousandeyes-cli) for API access. It selects tests, waits for agent results, evaluates thresholds, and produces JSON/JUnit evidence. The thin GitHub Action adds pipeline outputs and a job summary.

## Choose your platform

Start with the setup guide for your platform. Each integration runs `te-synthetics run-tests` with the same version 1 suite.

| Platform | Setup guide | Example |
| --- | --- | --- |
| GitHub Actions | [Action setup, inputs, and outputs](docs/platforms.md#github-actions) | [Workflow](examples/github-actions/ci.yml) |
| Developer terminal | [Native command](docs/platforms.md#developer-terminal) | [Suite](examples/synthetics/thousandeyes-ci.json) |
| Jenkins | [Worker and credential setup](docs/platforms.md#jenkins) | [Jenkinsfile](examples/jenkins/Jenkinsfile) |
| GitLab CI | [Runner and variable setup](docs/platforms.md#gitlab-ci) | [Pipeline](examples/gitlab/ci.yml) |
| Container | [Build and run an image](docs/platforms.md#container) | [Dockerfile](Dockerfile) |
| Argo CD | [PostSync validation](docs/platforms.md#argo-cd) | [PostSync Job](examples/argocd/post-sync-job.yaml) |
| Argo Rollouts | [Progressive delivery analysis](docs/platforms.md#argo-rollouts) | [AnalysisTemplate](examples/synthetics/argo-rollouts/analysis-template.yaml) and [rollout step](examples/synthetics/argo-rollouts/rollout-step.yaml) |

Cloud or Enterprise Agents perform the measurements. Select Enterprise Agents with connectivity to private targets. Each run creates new instant tests and consumes applicable ThousandEyes units.

## Configure your checks

Copy the [example suite](examples/synthetics/thousandeyes-ci.json), then replace its account group, agents, targets, and thresholds. Supply `TE_TOKEN` through your platform's secret store or shell environment.

Use the [JSON Schema](schemas/thousandeyes-ci.schema.json) for editor completion and generation. The example references it through `$schema`. See [suite configuration and results](docs/configuration.md) for threshold units, execution rules, command options, exit codes, and report semantics.

## Installation and version pinning

Use the [source build instructions](DEVELOPMENT.md#build-the-native-bundle) to obtain the runner and its compatible backend. Keep both executables together. [Backend metadata](backend.json) pins the official [ThousandEyes CLI 0.0.4 release](https://github.com/thousandeyes/thousandeyes-cli/releases/tag/0.0.4) to an exact source commit.

Pin Action source to a reviewed full commit SHA. Selecting a published runner through `runner-version` is a separate pin; it requires an exact available release. While this repository is private, use a source-built bundle through `runner-path`. See [version pinning and release verification](RELEASING.md#consumer-version-pins) for details and the current distribution limits.

## Development and maintenance

See [DEVELOPMENT.md](DEVELOPMENT.md) for prerequisites, builds, coding conventions, and validation. [CONTRIBUTING.md](CONTRIBUTING.md) describes contribution policy and pull requests. [RELEASING.md](RELEASING.md) documents dependency updates, release checks, and publication.

See [architecture](docs/architecture.md) for component responsibilities and the official CLI backend contract.

## Community and support

Use [GitHub issues](https://github.com/thousandeyes/thousandeyes-synthetics-ci/issues) to report runner or Action bugs, ask usage questions, and suggest features. Check existing issues and pull requests before opening a new discussion. Remove tokens, internal hostnames, account details, and other sensitive information from examples and reports.

Participation follows our [code of conduct](CODE_OF_CONDUCT.md). Report suspected vulnerabilities through the private email channel in [SECURITY.md](SECURITY.md).

## License and attribution

Licensed under Apache-2.0. See [LICENSE](LICENSE), [NOTICE](NOTICE), and [licensing and third-party attribution](docs/licensing.md). Bundled dependencies retain their own licenses and notices.
