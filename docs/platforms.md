# Platform setup

Every integration uses a version 1 suite and a ThousandEyes API token from the platform's secret store. Copy the [example suite](../examples/synthetics/thousandeyes-ci.json) to your application repository as `thousandeyes-ci.json`, then replace its account group, agents, targets, and thresholds. For editor completion, copy the [JSON Schema](../schemas/thousandeyes-ci.schema.json) alongside it and update the suite's `$schema` path. See [suite configuration and results](configuration.md) for threshold units, defaults, and report semantics.

These examples create real instant tests and consume applicable units. Cloud or Enterprise Agents perform the measurements; private targets require an Enterprise Agent that can reach them. Protect suite configuration and reports according to the target and agent information they contain.

## Developer terminal

Build the native bundle using [DEVELOPMENT.md](../DEVELOPMENT.md#build-the-native-bundle), or install a reviewed bundle supplied by your maintainer. Keep `te-synthetics` and its pinned `thousandeyes` backend together in the same directory, and add that directory to your `PATH`. Windows bundles use `.exe` suffixes. Node is needed to build the bundle, not to run its executables.

Supply `TE_TOKEN` through your secret manager or shell environment, then run:

```bash
te-synthetics run-tests --config thousandeyes-ci.json \
  --json-report results.json --junit-report results.xml
```

The runner resolves its backend from `--api-cli-path`, then `TE_API_CLI_PATH`, then the sibling `thousandeyes` executable, then `PATH`. An override must support the API options and response contract documented in [architecture](architecture.md#api-backend-contract).

The command's exit status controls the gate on every platform. A separate JUnit gate must account for non-blocking tests: their assertions remain XML failures even when the runner exits successfully. See [results and exit status](configuration.md#results-and-exit-status).

## GitHub Actions

Copy the [GitHub Actions workflow](../examples/github-actions/ci.yml) into your application's `.github/workflows/` directory, and commit the configured suite as `thousandeyes-ci.json`. Add a repository or environment secret named `TE_TOKEN`. The example uses `workflow_dispatch` so you can choose when to create instant tests. The workflow file must be on the application's default branch to appear in the manual dispatch interface.

The workflow checks out one full, reviewed source commit, builds its native bundle, and invokes the Action from that same checkout using `uses: ./.te-synthetics`. Update `TE_SYNTHETICS_REF` only after reviewing the replacement commit. The other Actions are also pinned to full commit SHAs. Pinning the wrapper alone does not pin a separately installed runner or its backend; keep those choices aligned. See [release and dependency guidance](../RELEASING.md) before updating pins.

If this project is private and the workflow runs in a different repository, create `TE_SYNTHETICS_READ_TOKEN` with **Contents: read** access to this source repository and the necessary organization authorization. The application's `GITHUB_TOKEN` is scoped to its own repository and cannot check out another private repository. The example passes the dedicated token only to checkout and disables persisted Git credentials. Once the source is public, or when running within this repository, the default `github.token` is sufficient. The ThousandEyes API token is supplied only to the synthetic-check Action step, after the build.

The source-build example uses `runner-path` and does not depend on a published release. Alternatively, when a public runner release is available, set `runner-version` to its exact version and pin the Action's remote `uses:` reference to a full reviewed commit SHA. The Action downloads the two-binary native bundle and verifies its archive SHA-256 against the release checksum manifest. Floating `latest` is rejected. The release downloader does not provide authentication for private release assets; use a source build or a separately managed bundle while the repository is private.

### Inputs

| Input | Default | Purpose |
| --- | --- | --- |
| `api-token` | Required | ThousandEyes API bearer token from a secret. |
| `config-path` | `thousandeyes-ci.json` | Suite file relative to the workspace, or absolute. |
| `runner-path` | Empty | Absolute path to a managed `te-synthetics` executable. |
| `runner-version` | Empty | Exact published release from this repository; required if no runner path is supplied. |
| `api-cli-path` | Bundled backend | Absolute path to a compatible official API CLI override. |
| `account-group-id` | Suite value | Override the account group. |
| `timeout`, `poll-interval` | Suite values | Positive duration overrides such as `5m` and `5s`. |
| `api-base-url` | `https://api.thousandeyes.com` | API origin; HTTP is permitted only for loopback test servers. |

`cli-version` and `cli-path` are deprecated aliases for `runner-version` and `runner-path`. Conflicting aliases are rejected. [action.yml](../action.yml) defines the complete input and output contract.

### Outputs and reports

The Action exposes `run-id`, `exit-code`, `json-report`, `junit-report`, and test counts (`passed`, `failed`, `non-blocking-failed`, `skipped`, `timed-out`, `errors`). Report paths are reserved before execution; files exist only if the runner produces them. The example uploads reports with `if: always()` so a failed check can retain its evidence. Choose artifact access and retention settings appropriate to the target and agent information in the reports.

The job summary includes timestamps, duration, test IDs, status symbols, and agent tables with test name, target, agent name/ID, round, measurements, bounded failure reasons, and permitted ThousandEyes result links. Failure rows are prioritized when summary limits are reached. Complete JSON/JUnit artifacts retain all results. The Action validates reports and preserves the runner's exit decision.

The token is masked and passed through `TE_TOKEN`, never command arguments. Raw child output is not forwarded to GitHub logs. Summary values are bounded, sanitized, token-redacted, and escaped. Run secret-bearing workflows only from trusted code and refs; do not add an untrusted checkout to a privileged workflow such as `pull_request_target`.

## Jenkins

Install the [native bundle](../DEVELOPMENT.md#build-the-native-bundle) on the worker's `PATH`, keeping both executables together. Adapt the [Jenkinsfile](../examples/jenkins/Jenkinsfile) and create a secret-text credential named `thousandeyes-token`.

The pipeline uses the command's exit status for the gate and archives JSON/JUnit files after a failed check. If you add an independent JUnit gate, account for [non-blocking assertions](configuration.md#results-and-exit-status). Validate credential access, worker connectivity, and artifact permissions in your Jenkins environment.

## GitLab CI

Install the [native bundle](../DEVELOPMENT.md#build-the-native-bundle) on the worker's `PATH`, keeping both executables together. Copy [ci.yml](../examples/gitlab/ci.yml) to `.gitlab-ci.yml`, configure a masked/protected `TE_TOKEN`, and select an appropriately tagged Linux shell runner.

The example runs only from the web UI. Protected variables require an eligible protected ref. It uses the command's exit status for the gate and archives JSON/JUnit files after a failed check. Validate runner permissions, protected-variable availability, and artifact access in your GitLab environment.

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

Run these commands from a reviewed source checkout with Node.js and Docker installed. Supply `TE_TOKEN` through your secret manager or shell environment before running the container. The container entrypoint is `te-synthetics`. Both native executables are included, so no Node, shell, or writable temporary directory is required at runtime. `--json` writes a sanitized report to stdout. For report files, mount a writable output volume and pass its paths to `--json-report` and `--junit-report`.

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

Complete the [Kubernetes prerequisites](#kubernetes-prerequisites), then adapt the [PostSync Job](../examples/argocd/post-sync-job.yaml). Argo CD runs the example after synchronization and treats a failed Job as a failed hook. Validate the hook's namespace, target revision, secret access, and image digest in your Argo CD environment.

## Argo Rollouts

Complete the [Kubernetes prerequisites](#kubernetes-prerequisites), apply the [AnalysisTemplate](../examples/synthetics/argo-rollouts/analysis-template.yaml), and add the [rollout step](../examples/synthetics/argo-rollouts/rollout-step.yaml). Supply the image, account group, suite ConfigMap, and Secret arguments. Argo Rollouts uses the Job's completion to decide whether the analysis succeeded. A ClusterAnalysisTemplate can use the same spec when appropriate. Validate the analysis against the intended preview or canary revision before using it to gate a rollout.

Jenkins, GitLab, and Argo runtime examples require validation in your own deployment environment; they are not exercised end to end by this repository's CI.
