# Releases

## Consumer version pins

| Component | Use |
| --- | --- |
| GitHub Action | A reviewed full commit SHA. |
| Runner | The version pinned by the Action, a reviewed source build, or an exact published `runner-version` override without a `v` prefix. |
| Official CLI backend | The source commit in [backend.json](backend.json). |
| Container | An image digest. |

The [GitHub Actions example](examples/github-actions/ci.yml) builds the runner and loads the Action from the same pinned checkout. Published Action releases with a runner pin download the matching executable. Action version `1.0.0` requires an explicit version or executable path.

While the repository is private, use a source build with `runner-path`; release downloads do not support GitHub authentication. Archive SHA-256 checksums detect mismatches with the release manifest, but do not establish publisher identity.

## Update dependencies

Review changes to `backend.json`, `package-lock.json`, `go.mod`, `go.sum`, and the backend's dependencies. Rebuild `dist/` after Action dependency changes, run the [development checks](DEVELOPMENT.md#validate-changes), and review [required license notices](docs/licensing.md). Leave `TE_API_CLI_SOURCE` unset for releases.

## Publish a release

1. Choose an unused `X.Y.Z` version, or a prerelease version when using a preview backend. Set `DEFAULT_RUNNER_VERSION` in [src/install.js](src/install.js) to this version and rebuild and commit `dist/`.
2. Select a reviewed `main` commit with the updated pin and passing CI. Complete the required publishing approvals. Pushing a version tag starts publication immediately. Set `VERSION` and `REVIEWED_SHA`, then run:

   ```bash
   git tag -a "$VERSION" "$REVIEWED_SHA" -m "Release $VERSION"
   git push origin "refs/tags/$VERSION"
   ```

3. Check the [release workflow](.github/workflows/release.yml). It tests and publishes native bundles and checksums, then the container image. GitHub releases are currently marked prerelease.
4. Verify installation, the runner version, and consumer access. Record the Action SHA, runner version, backend commit, and image digest before announcing the release.

Native publication can succeed before a container push fails. Check both results before retrying. Correct a published version with a new release rather than moving its tag or replacing assets.
