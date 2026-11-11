# Releases and version pinning

## Consumer version pins

| Component | Pin |
| --- | --- |
| GitHub Action | A reviewed full 40-character commit SHA. |
| Native runner | A reviewed source commit, or an exact published `runner-version`. |
| Official CLI backend | The full source commit in [backend.json](backend.json). |
| Container | The image digest from your registry. |

Action and runner pins are independent. The [GitHub Actions example](examples/github-actions/ci.yml) builds the runner and loads the Action from one pinned checkout. Other Actions also use full commit SHAs; see [GitHub's secure use guidance](https://docs.github.com/en/actions/reference/security/secure-use#using-third-party-actions).

`runner-version` accepts exact `X.Y.Z` or prerelease versions without a `v` prefix. `latest`, ranges, and build metadata are unsupported. Choose an available [release](https://github.com/thousandeyes/thousandeyes-synthetics-ci/releases); the `package.json` version does not establish one, and npm publication is disabled.

Release downloads do not support GitHub authentication. While this repository is private, build from an authenticated checkout and use `runner-path`. The ThousandEyes `api-token` does not authenticate GitHub downloads.

The installer checks archive SHA-256 values against the manifest from the same release. These checksums detect mismatches; they are not signatures or provenance attestations.

## Dependencies and distribution

Review changes to [backend.json](backend.json), [package-lock.json](package-lock.json), [go.mod](go.mod), [go.sum](go.sum), and the backend's own dependencies. Rebuild `dist/` after Action dependency changes and rerun the [development checks](DEVELOPMENT.md#validate-changes). Leave `TE_API_CLI_SOURCE` unset for releases.

Archives contain both executables, the root `LICENSE`, `NOTICE`, `README.md`, and `backend.json`. The image includes both executables and the root license, notice, and backend metadata. Packaging does not automatically collect every third-party notice; review and include required texts as described in [licensing](docs/licensing.md).

The GoReleaser tool version range, hosted runner images, and Docker base tags can change. Source pins alone do not make the complete build reproducible.

## Maintainer release procedure

Pushing a version tag starts the [release workflow](.github/workflows/release.yml) immediately. Required approvals and selection of a reviewed `main` commit are maintainer responsibilities; the workflow does not enforce those gates.

1. Select a reviewed `main` commit with all [CI checks](.github/workflows/ci.yml) passing. Complete required security, legal, licensing, and publishing approvals, and verify dependency pins, generated files, and distribution notices.
2. Choose an unused exact version without a `v` prefix; preview backends require a prerelease suffix. Set `VERSION` and `REVIEWED_SHA`, then tag the reviewed commit:

   ```bash
   git tag -a "$VERSION" "$REVIEWED_SHA" -m "Release $VERSION"
   git push origin "refs/tags/$VERSION"
   ```

3. Check the workflow results. Tests, bundle verification, and container smoke precede publication. The workflow creates a GitHub prerelease with six native archives (Linux, macOS, and Windows on amd64/arm64) and checksums, then publishes the Linux amd64/arm64 image to `ghcr.io/thousandeyes/thousandeyes-synthetics-ci:<version>`.
4. Verify the installed runner version, both publication results, and consumer access. Record the Action SHA, runner version, backend commit, and image digest before announcing the release.

All GitHub releases are currently marked prerelease, including plain `X.Y.Z` tags. Native publication can succeed before a container push fails. Inspect existing assets before retrying; correct published versions with a new release rather than moving a tag or replacing assets.
