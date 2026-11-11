# Development

This guide contains development setup, build instructions, coding conventions, and local checks. See [CONTRIBUTING.md](CONTRIBUTING.md) for contribution policy and [RELEASING.md](RELEASING.md) for versioning and publication.

## Prerequisites

Install:

- Git.
- The Go version specified in [go.mod](go.mod), or a compatible newer version.
- Node.js 24 and npm. The Action runs on Node.js 24 in GitHub Actions.
- **Optional:** Docker to check the runtime image. The container API smoke test requires Linux.

## Build the native bundle

Clone this repository, or your fork when contributing, and run:

```bash
git clone https://github.com/thousandeyes/thousandeyes-synthetics-ci.git
cd thousandeyes-synthetics-ci
node scripts/build-native.js
```

Access to a private repository requires an authenticated Git client. To build a specific reviewed revision, check out its full commit SHA before running the build.

The native build downloads the official CLI source pinned in [backend.json](backend.json) and Go dependencies. It compiles both executables into `build/bin`: `te-synthetics` and `thousandeyes`, with `.exe` suffixes on Windows. Keep them together when copying the bundle. Node is a build dependency; the native binaries run without Node.

See [platform setup](docs/platforms.md#developer-terminal) for backend discovery and [architecture](docs/architecture.md#api-backend-contract) for the required API contract.

For local backend development, `TE_API_CLI_SOURCE` can point to an official CLI checkout. Preparation copies source-control-visible files, including untracked additions that are not ignored. This override bypasses the remote source pin; leave it unset for release builds.

For suite setup and execution, follow the [platform guide](docs/platforms.md). Building the bundle does not execute live ThousandEyes tests.

## Work on the GitHub Action

Install the locked dependencies without lifecycle scripts:

```bash
npm ci --ignore-scripts
```

Edit JavaScript under `src/`, then run `npm run build` to regenerate the Action entry point and third-party license texts. Commit `dist/` with the corresponding source or dependency changes; do not edit it directly. Release artifacts go under `build/release`, separate from the committed Action bundle.

## Coding and documentation conventions

- Format Go changes with `gofmt` or `make fmt`.
- Follow the existing CommonJS style in JavaScript.
- Keep API access in the official CLI backend and CI policy in the runner. Keep the Action focused on integration and reporting.
- Preserve compatibility for suite and report formats, or describe a proposed breaking change explicitly.
- Use inclusive language, descriptive names, and concise documentation. Examples must use synthetic targets, identities, and credentials.
- Preserve existing copyright notices, license headers, and third-party attribution. Add a header to new source files as described in [licensing](docs/licensing.md#source-headers).
- Keep user setup in [platforms](docs/platforms.md), the suite contract in [configuration](docs/configuration.md), development instructions here, and release operations in [RELEASING.md](RELEASING.md). Link to the relevant guide instead of duplicating commands in the README.

## Validate changes

Run checks that cover the behavior you changed. After installing npm dependencies and building the native bundle, run:

```bash
npm test
npm run build
TE_TEST_API_CLI_PATH="$PWD/build/bin/thousandeyes" make ci
npm run smoke:native
ACTION_ENTRY=dist/index.js node --test test/action.test.js
```

These checks use fake tokens and a local fake API, so no real ThousandEyes account is required. Native tests run the actual pinned official CLI. On Windows, use `thousandeyes.exe` for the backend test path and set environment variables with your shell's syntax. Run `go test ./...` and `go vet ./...` directly if Make is unavailable.

Native command integration tests explicitly skip when `TE_TEST_API_CLI_PATH` is absent. Skipped integration tests do not validate the backend contract. After committing generated Action changes, `npm run check:dist` verifies that rebuilding produces no differences.

For container changes, run these additional checks on Linux with Docker available:

```bash
node scripts/prepare-backend.js
node scripts/smoke-container.js --api
```

The [CI workflow](.github/workflows/ci.yml) tests Linux, macOS, and Windows, a restricted container, and six-platform release snapshots. Platform examples require separate validation in the relevant CI or Kubernetes environment.

Run live tests only against targets and agents you are authorized to use. Instant tests consume applicable ThousandEyes units. Do not commit real credentials or reports containing private account data.
