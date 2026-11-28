# Platform setup

Copy the [example suite](../examples/synthetics/thousandeyes-ci.json) to your application repository as `thousandeyes-ci.json`, then replace its account group, agents, targets, and thresholds. Supply a ThousandEyes API token through your platform's secret store. For editor completion, copy the [JSON Schema](../schemas/thousandeyes-ci.schema.json) and update the suite's `$schema` path. See [configuration and results](configuration.md) for the shared contract.

Instant tests consume applicable units. Private targets require Enterprise Agents that can reach them. Review [report access and retention](#report-access-and-retention) before running checks with private data.

## Report access and retention

Reports intentionally retain test/account/agent identifiers, names, targets, measurements, and result links. Token redaction and URL sanitization do not anonymize this metadata. Use an appropriately restricted private project for confidential live runs. Public repositories should use fake test data or data explicitly approved for public sharing.

GitHub artifacts are downloadable by [signed-in repository readers](https://docs.github.com/en/actions/how-tos/manage-workflow-runs/download-workflow-artifacts). The Action also writes metadata to the job summary, and CLI output can enter logs. Removing the artifact upload does not hide the summary. Review the repository's audience before running checks.

| Platform | Access and cleanup |
| --- | --- |
| GitHub Actions | The example and repository CI use `retention-days: 7`. Check the uploaded artifact's actual expiry and applicable repository/organization limits. See [artifact retention](https://github.com/actions/upload-artifact#retention-period). |
| GitLab CI | The example uses `expire_in: 7 days`. Review project/job access and [keep-latest-successful settings](https://docs.gitlab.com/ci/jobs/job_artifacts/#keep-artifacts-from-most-recent-successful-jobs), which can retain artifacts beyond that period. |
| Jenkins | Restrict job access and configure build/log/artifact retention and workspace cleanup. Archiving a report alone does not configure deletion. |
| Argo / Kubernetes | The examples print JSON to pod logs. Restrict log access and configure Job cleanup and [centralized log retention](https://kubernetes.io/docs/concepts/cluster-administration/logging/) separately. |
| Developer terminal | Write reports to restricted paths. The owner controls deletion of files and any captured terminal output. |

Seven days is a sample artifact lifetime, not a Cisco retention policy. Artifact expiry does not shorten retention for logs, summaries, workspaces, committed suites, or exported copies.

Before using an integration with private data:

1. Run a job with non-confidential data; this repository's [CI](../.github/workflows/ci.yml) uses a local mock API. Record its run ID or URL and reviewed commit SHA.
2. Check artifact and run-page access for intended readers, including denied access with a test account. Verify the created artifact's expiry, rather than relying only on workflow configuration.
3. Record the separate log/summary retention and workspace cleanup settings, including any keep-latest or manual-retention exceptions. Repeat the checks when permissions, repository visibility, or retention settings change.

## Developer terminal

Build the [native bundle](../DEVELOPMENT.md#build-the-native-bundle), or install a reviewed bundle supplied by your maintainer. Keep both executables together and add their directory to `PATH`. Windows bundles use `.exe` suffixes.

Supply `TE_TOKEN` through your secret manager or shell environment, then run:

```bash
te-synthetics run-tests --config thousandeyes-ci.json \
  --json-report results.json --junit-report results.xml
```

The runner resolves its backend from `--api-cli-path`, then `TE_API_CLI_PATH`, then the sibling `thousandeyes` executable, then `PATH`. An override must support the API options and response contract documented in [architecture](architecture.md#api-backend-contract).

The command's exit status controls the gate on every platform. A separate JUnit gate must account for non-blocking tests: their assertions remain XML failures even when the runner exits successfully. See [results and exit status](configuration.md#results-and-exit-status).

## GitHub Actions

Copy the [workflow](../examples/github-actions/ci.yml) into your application's `.github/workflows/` directory and commit the configured suite as `thousandeyes-ci.json`. Add a repository secret named `TE_TOKEN`. The example runs manually through `workflow_dispatch`; put it on the default branch to make it available in the Actions interface.

The workflow builds the runner and loads the local Action from one pinned checkout. Update `TE_SYNTHETICS_REF` after reviewing the replacement commit. See [version pinning](../RELEASING.md#consumer-version-pins) for separately installed runners and releases.

To check out this private source from another repository, add `TE_SYNTHETICS_READ_TOKEN` with **Contents: read** access and any required organization authorization. The application's `GITHUB_TOKEN` cannot read another private repository. Public or same-repository checkouts use the default token. Checkout disables persisted credentials; the ThousandEyes token is supplied only to the test step.

The source-build example uses `runner-path`. When you use a published Action release that pins its runner, the Action downloads the matching bundle automatically. Action version `1.0.0` requires an explicit `runner-version` or `runner-path`. See [release guidance](../RELEASING.md#consumer-version-pins) for version overrides and private-download limits.

### Inputs

| Input | Default | Purpose |
| --- | --- | --- |
| `api-token` | Required | ThousandEyes API bearer token from a secret. |
| `config-path` | `thousandeyes-ci.json` | Suite file relative to the workspace, or absolute. |
| `runner-path` | Empty | Absolute path to a managed `te-synthetics` executable. |
| `runner-version` | Version pinned by the Action | Override the runner download with an exact published release. |
| `api-cli-path` | Bundled backend | Absolute path to a compatible official API CLI override. |
| `account-group-id` | Suite value | Override the account group. |
| `timeout`, `poll-interval` | Suite values | Positive duration overrides such as `5m` and `5s`. |
| `api-base-url` | `https://api.thousandeyes.com` | API origin; HTTP is permitted only for loopback test servers. |

`cli-version` and `cli-path` are deprecated aliases for `runner-version` and `runner-path`. Conflicting aliases are rejected. [action.yml](../action.yml) defines the complete input and output contract.

### Outputs and reports

The Action exposes `run-id`, `exit-code`, `json-report`, `junit-report`, and test counts (`passed`, `failed`, `non-blocking-failed`, `skipped`, `timed-out`, `errors`). Report paths are reserved before execution; files exist only if the runner produces them. The example uses `if: always()` to upload reports after failed checks.

The job summary shows test and agent results, measurements, failure reasons, and ThousandEyes links. It prioritizes failures when space is limited; JSON/JUnit artifacts retain all results. The runner's exit decision controls the gate.

The token is masked and passed through `TE_TOKEN`, never command arguments. Raw child output is not forwarded to logs, and summary values are bounded, sanitized, token-redacted, and escaped.

## Jenkins

Follow the [terminal setup](#developer-terminal) on the worker. Adapt the [Jenkinsfile](../examples/jenkins/Jenkinsfile) and create a secret-text credential named `thousandeyes-token`. The pipeline archives JSON/JUnit reports after failed checks.

## GitLab CI

Follow the [terminal setup](#developer-terminal) on a Linux shell runner. Copy [ci.yml](../examples/gitlab/ci.yml) to `.gitlab-ci.yml`, configure a masked/protected `TE_TOKEN`, and match the runner tags. The example runs from the web UI and archives reports after failed checks. Protected variables require an eligible protected ref.

## Container

Build the pinned backend source, then the image:

```bash
node scripts/prepare-backend.js
docker build -t te-synthetics:local .
docker run --rm --read-only --cap-drop=ALL \
  --security-opt=no-new-privileges -e TE_TOKEN \
  -v "$PWD/examples/synthetics:/config:ro" \
  te-synthetics:local run-tests \
  --config /config/thousandeyes-ci.json --json
```

Run from a reviewed checkout with Node.js and Docker installed, and supply `TE_TOKEN` in the environment. The image includes both executables and uses `te-synthetics` as its entrypoint. `--json` writes the report to stdout. To save files, mount a writable output volume and set `--json-report` and `--junit-report`.

Release images stamp the release tag into `te-synthetics --version`. Local builds report `dev` unless you pass `--build-arg RUNNER_VERSION=<version>` to `docker build`. See [release guidance](../RELEASING.md) for image publishing and dependency review.

### Kubernetes prerequisites

Build and publish the image to your own registry, then pin its digest in your Argo configuration. Create a Secret named `thousandeyes-api` with a `token` key, and create the suite ConfigMap:

```bash
kubectl create configmap thousandeyes-canary-suite \
  --from-file=thousandeyes-ci.json=examples/synthetics/thousandeyes-ci.json
```

The Job, ConfigMap, and Secret must be in the same namespace. The examples disable Job retries to avoid creating duplicate instant tests and set the outer Job deadline longer than the runner's five-minute deadline. Keep completed Job logs or export reports to your evidence store; these examples do not add external artifact storage.

Configure the suite target to reach the canary or preview revision explicitly. A weighted shared endpoint can route to the stable revision. A Job running inside Kubernetes does not give remote agents access to cluster-internal Services. Select Enterprise Agents with connectivity to private targets.

## Argo CD

Complete the [Kubernetes prerequisites](#kubernetes-prerequisites), then adapt the [PostSync Job](../examples/argocd/post-sync-job.yaml). Argo CD runs it after synchronization and treats a failed Job as a failed hook.

## Argo Rollouts

Complete the [Kubernetes prerequisites](#kubernetes-prerequisites), apply the [AnalysisTemplate](../examples/synthetics/argo-rollouts/analysis-template.yaml), and add the [rollout step](../examples/synthetics/argo-rollouts/rollout-step.yaml). Supply the image, account group, suite ConfigMap, and Secret arguments. Argo Rollouts uses the Job's completion to decide whether the analysis succeeded.

Jenkins, GitLab, and Argo runtime examples require validation in your own deployment environment; they are not exercised end to end by this repository's CI.
