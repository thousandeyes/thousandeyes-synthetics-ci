# Development

See [CONTRIBUTING.md](CONTRIBUTING.md) for contribution policy and [RELEASING.md](RELEASING.md) for publication.

## Prerequisites

Install Git, the Go version in [go.mod](go.mod), and Node.js 24 with npm. Container checks also require Docker on Linux.

## Build the native bundle

Clone this repository, or your fork, and run:

```bash
git clone https://github.com/thousandeyes/thousandeyes-synthetics-ci.git
cd thousandeyes-synthetics-ci
node scripts/build-native.js
```

Private source requires an authenticated Git client. Check out a reviewed full commit SHA before building a pinned version.

The build compiles `te-synthetics` and the official `thousandeyes` backend pinned in [backend.json](backend.json) into `build/bin`, with `.exe` suffixes on Windows. Keep both executables together. Node is needed for building, not for running them. See [platform setup](docs/platforms.md) for usage.

For local backend development, set `TE_API_CLI_SOURCE` to an official CLI checkout. This bypasses the remote pin; leave it unset for releases. See [architecture](docs/architecture.md#api-backend-contract) for the backend contract.

## Coding conventions

- Format Go with `gofmt` or `make fmt`; follow existing CommonJS style in JavaScript.
- Keep API access in the official CLI, CI policy in the runner, and integration/reporting in the Action.
- Preserve suite and report compatibility, or explain proposed breaking changes.
- Use descriptive names, concise documentation, and synthetic values in examples.
- Preserve copyright and third-party notices; follow [source-header guidance](docs/licensing.md#source-headers) for new files.
- Edit Action source under `src/`; rebuild and commit `dist/` with related source or dependency changes. Do not edit generated files directly.
- Link to platform, configuration, development, and release guides instead of repeating their instructions in the README.

## Validate changes

After building the native bundle, run the checks relevant to your changes:

```bash
npm ci --ignore-scripts
npm test
npm run build
TE_TEST_API_CLI_PATH="$PWD/build/bin/thousandeyes" make ci
npm run smoke:native
ACTION_ENTRY=dist/index.js node --test test/action.test.js
```

Checks use fake tokens and a local API. Native integration tests run the actual pinned CLI and skip when `TE_TEST_API_CLI_PATH` is absent. On Windows, use `thousandeyes.exe` and your shell's environment-variable syntax. If Make is unavailable, run `go test ./...` and `go vet ./...` directly.

Commit rebuilt `dist/` when needed, then run `npm run check:dist` to confirm it matches the source. For container changes, also run on Linux:

```bash
node scripts/prepare-backend.js
node scripts/smoke-container.js --api
```

[CI](.github/workflows/ci.yml) checks Linux, macOS, Windows, a restricted container, and six-platform release snapshots. Validate platform examples in their deployment environment. Run live tests only against authorized targets and agents; do not commit credentials or private reports.
