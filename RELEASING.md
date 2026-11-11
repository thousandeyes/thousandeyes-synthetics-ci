# Releases and version pinning

This guide covers version selection for users and publication by maintainers. See [DEVELOPMENT.md](DEVELOPMENT.md) for build and test commands.

## Consumer version pins

The Action, runner, API backend, and container have separate version selections:

| Component | What to pin | How to update |
| --- | --- | --- |
| GitHub Action code | A reviewed full 40-character commit SHA. | Review changes to `action.yml`, `src/`, `dist/`, and dependencies before changing the SHA. |
| Native runner | Build from a reviewed source commit, or set `runner-version` to an exact published version such as `0.1.0` if that release exists. | Review the selected source or release, including its backend metadata. |
| Official CLI backend | The full source commit in [backend.json](backend.json). | Review the official CLI release and commit, then rebuild and test the bundle. |
| Container | An image digest from your registry. | Validate the new image and update the digest in your deployment or analysis configuration. |

Pinning Action code does not pin a runner downloaded through `runner-version`. To keep source selections together, the [GitHub Actions example](examples/github-actions/ci.yml) checks out one full source SHA, builds the runner there, and runs the Action from that local checkout. Other third-party Actions in the example also use full commit SHAs. See [GitHub's secure use guidance](https://docs.github.com/en/actions/reference/security/secure-use#using-third-party-actions).

The installer accepts exact `X.Y.Z` or prerelease versions, without a `v` prefix. It rejects `latest`, version ranges, and build metadata. A version in `package.json` does not establish an available runner release; npm publication is disabled with `private: true`. Check the repository's [releases](https://github.com/thousandeyes/thousandeyes-synthetics-ci/releases) before choosing a version.

The installer fetches archive and checksum URLs without a GitHub credential. While this repository is private, build from an authenticated source checkout and use `runner-path`. The ThousandEyes `api-token` input authenticates ThousandEyes API requests, not GitHub downloads.

For accessible published releases, the installer verifies the archive's SHA-256 against the checksum manifest downloaded from the same release. This detects a mismatch between those files; it does not independently prove their source or publisher. The current release workflow does not generate artifact signatures or provenance attestations. Use an independently verified, managed bundle through `runner-path` when your environment requires that assurance.

## Dependencies and distribution contents

- [backend.json](backend.json) pins the official CLI repository, release, and source commit. Preparation verifies the fetched commit. Leave `TE_API_CLI_SOURCE` unset when building releases; it is a development override.
- [package-lock.json](package-lock.json) records Action dependencies and integrity values. Use `npm ci --ignore-scripts`, review dependency changes, and commit regenerated `dist/index.js` and `dist/licenses.txt` together.
- [go.mod](go.mod) and [go.sum](go.sum) record runner dependencies. The prepared backend has its own Go module files. Review dependencies for both executables.
- Workflow Actions use commit SHAs, but the GoReleaser tool constraint, hosted runner images, and Docker base image tags can change. Source pinning alone does not make the complete build reproducible.

Native archives contain both executables, the root `LICENSE`, `NOTICE`, `README.md`, and `backend.json`. Linux and macOS archives use `.tar.gz`; Windows archives use `.zip`. Each OS has amd64 and arm64 builds. The runtime image contains both Linux executables and the root license, notice, and backend metadata.

Before distributing either format, review and include applicable backend, dependency, and base-image license texts and notices. The packaging configuration currently copies the root license and notice; it does not automatically collect every third-party notice. See [licensing and attribution](docs/licensing.md).

## Maintainer release procedure

Pushing a matching remote version tag immediately starts the [release workflow](.github/workflows/release.yml). Complete the required security, legal, licensing, and publishing approvals before pushing a release tag. The workflow itself has no approval environment or check that the tag points to an approved `main` commit.

1. Choose a reviewed commit on `main` whose [CI](.github/workflows/ci.yml) has passed Linux, macOS, Windows, container, and release-snapshot checks. Confirm that generated Action files match their sources and that the distribution notices are complete.
2. Verify the selected official CLI release and full source commit in `backend.json`. Rerun integration checks after any backend or dependency update. Confirm that `TE_API_CLI_SOURCE` is unset in the release environment.
3. Choose an unused exact semantic version without a `v` prefix. Use a prerelease suffix when required by a preview backend. The workflow currently marks every GitHub release as a prerelease, including plain `X.Y.Z` tags.
4. After publication is authorized, create the tag at the reviewed full commit SHA and push that tag. Use explicit values for `VERSION` and `REVIEWED_SHA`:

   ```bash
   git tag -a "$VERSION" "$REVIEWED_SHA" -m "Release $VERSION"
   git push origin "refs/tags/$VERSION"
   ```

5. Monitor the release workflow. It runs Node tests, verifies the generated Action bundle, builds the native distribution, runs Go integration and container smoke tests, and builds and verifies all six archives before creating the GitHub release. Archive verification uses the Action installer and tests native startup on the build host. It does not execute every target architecture on that host.
6. Confirm that the GitHub prerelease contains six native archives and its checksum manifest. The workflow uses `gh release create --verify-tag`; that option confirms the remote tag exists, rather than verifying a tag signature.
7. Confirm that the separate container publication step succeeds for Linux amd64 and arm64 at `ghcr.io/thousandeyes/thousandeyes-synthetics-ci:<version>`. Both native and container release builds stamp the runner version. Record the image digest and verify package visibility and consumer access before announcing it.
8. Validate the installed bundle or image, its reported runner version, and the intended consumer workflow. Document the Action commit SHA, runner version, backend commit, and image digest relevant to that release.

The native release is created before the container push. A later image-publishing failure can leave a native release without its corresponding container. Inspect the existing assets and workflow state before retrying; do not move an announced tag or silently replace a published version. Correct a published release with a new version.
