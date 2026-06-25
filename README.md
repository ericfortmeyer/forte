<!-- markdownlint-disable MD033-->
<h1 align="center">Forte</h1>
<p align="center">
  Focused deploy tool that enforces FHS conventions so deployments are predictable.
</p>

<p align="center">
  <img src="https://github.com/ericfortmeyer/forte/actions/workflows/push.yml/badge.svg" alt="Quality Checks">
  <img src="https://github.com/ericfortmeyer/forte/actions/workflows/release.yml/badge.svg?event=release" alt="Release">
  <img src="https://github.com/ericfortmeyer/forte/actions/workflows/pr.yml/badge.svg" alt="Integration Tests">
  <a href='https://coveralls.io/github/ericfortmeyer/forte'><img src='https://coveralls.io/repos/github/ericfortmeyer/forte/badge.svg' alt='Coverage Status' /></a>
</p>

## Why Forte

Developers and operators often disagree where files should live: examples scatter across /var/www, /srv, distro-specific vhost paths, and repos sometimes use git as an ad-hoc deploy mechanism. Forte removes that guesswork by enforcing simple, FHS-friendly defaults: runtime assets under `/srv/<app>` and system configuration under `/etc/<app>`. That makes deployments predictable, simplifies ownership and backup policies, and fits existing tooling (service units, SELinux, packaging) so teams spend less time wiring deployment logic and more time shipping code.

## Quick start

Download a binary from <https://github.com/ericfortmeyer/forte/releases>

**Note:** Not available on Windows

Run the CLI:

```zsh
./forte help
```

Deploy (example):

```zsh
sudo ./forte deploy myapp serviceuser
```

Result:

- Application binaries and runtime files → /srv/myapp/
- Config files → /etc/myapp/
- **Static assets** → /srv/assets/myapp/

To deploy assets, place them in `/tmp/myapp-assets` before running `deploy`. `/tmp/myapp-assets.tar.gz` is also a supported deployment source.

## Minimal API

Usage:

```zsh
./bin/forte help
usage: forte <command> [<args>]

  forte help                          Show this help
  forte version                       Display Forte version
  forte deploy <app-name> <user-name> Deploy an application
```

Behavior summary:

- No configuration or flags.
- Expected source layout: `/tmp/<app>`, `/tmp/<app>-config`, and `/tmp/<app>-assets` or `/tmp/<app>.tar.gz`, `/tmp/<app>-config.tar.gz`, and `/tmp/<app>-assets.tar.gz`.
- Mapping:
  - `/tmp/<app>` → `/srv/<app>`
  - `/tmp/<app>-config` → `/etc/<app>`
  - `/tmp/<app>-assets` → `/srv/assets/<app>/`
  - `/tmp/<app>.tar.gz` → `/srv/<app>`
  - `/tmp/<app>-config.tar.gz` → `/etc/<app>`
  - `/tmp/<app>-assets.tar.gz` → `/srv/assets/<app>/`
- The `-config` and `-assets` suffixes on source directories signal configuration and static assets, respectively.
- Tarball deployment sources is supported and permissions and ownership are preserved.
- Binaries released for amd64 and arm64 as GitHub artifacts.

## Guarantees and current limitations

- Idempotence: deployments are idempotent. Forte skips source files that are older than the corresponding destination files.
- Permissions: default ownership and modes applied on install:
  - Directories: root:serviceuser, 0750
  - Files:       root:serviceuser, 0640
- Known limitations:
  - Permission-only changes in the source are currently ignored if the file contents are unchanged. (Forte will not update permissions in that case.)
  - No dry-run option.
  - No automatic cleanup of the source directory.
  - No rollback of partial or failed deployments; partial state may remain on error.

## Troubleshooting notes (short)

- Forte expects source directories under /tmp. Ensure your CI or build step places the app at `/tmp/<app>` and config at `/tmp/<app>-config` before running deploy.
- If files are not updated, confirm the source file timestamp is newer than the destination. Permission-only changes are a known issue.
- If SELinux or permissions block the service, check contexts and ownership after deploy; fixes may be needed until permission-update behavior is addressed.

## Roadmap & investigations

Planned items and investigations to address current limitations:

- Fix permission-only sync behavior (decide whether to treat permission diffs as changes, or offer an explicit flag to sync perms).
- Implement safe rollback/transactional deploys or an atomic swap strategy to avoid partial state on failure.
- Optional: cleanup step (configurable) to remove or archive the source after successful deploy.
- Improve mapping flexibility (custom source/destination paths, non-/tmp sources).
- Publish package installers (deb/rpm) and integrate with distro packaging expectations.
- Add stronger CI enforcement (commit conventions, commitizen) and automated release artifacts for amd64/arm64.

## Short example

```zsh
# build pipeline places artifacts:
# /tmp/myApp /tmp/myApp-assets and /tmp/myApp-config

./bin/forte deploy myApp serviceuser
# -> /srv/myApp/ /srv/assets/myApp and /etc/myApp/
```

## Where Forte puts files — rationale and convention

Keep web application files grouped under a single, semantically meaningful tree: nest each app under `/srv/<app-name>`. This makes the layout explicit, predictable, and easy to explain to operators:

- **Clarity:** `/srv/<app-name>/` immediately signals “this is web infrastructure.”
- **Discoverability:** Tools, runbooks, and on-call engineers can find app files by server role instead of hunting across multiple top-level directories.
- **Separation of concerns:** Keep runtime assets belonging to the web app next to the app; keep global or shared asset stores distinct only when semantically necessary.

Avoid introducing a top-level /srv/assets/ that mixes different server roles. Splitting assets into a separate top-level path blurs the “what serves this” question and makes role-based policies (backup, permissions, SELinux contexts) harder to reason about.

### Asset deployment

Static assets (CSS, JavaScript, images, fonts, etc.) are deployed to `/srv/assets/<app>/` under a shared asset store, keeping them semantically distinct from application binaries while remaining part of the app's runtime deployment:

- **Source:** Place assets in `/tmp/<app>-assets/` or `/tmp/<app>-assets.tar.gz` before deploying.
- **Destination:** Assets deploy to `/srv/assets/<app>/`
- **Rationale:** A shared asset store simplifies CDN or reverse-proxy configurations, enables efficient cache invalidation per app, and keeps asset ownership clear.

**Example:**

Build pipeline prepares:

`/tmp/myapp` and `/tmp/myapp-assets` or `/tmp/myapp.tar.gz` and `/tmp/myapp-assets.tar.gz`

`./bin/forte deploy myapp serviceuser`
-> `/srv/myapp/`, `/etc/myapp/`, and `/srv/assets/myapp/`

Configure your web server to serve assets from `/srv/assets/<app>/` and app routes from `/srv/<app>/`

## FHS compliance

**Forte** follows the *Filesystem Hierarchy Standard* conventions: application runtime files live under `/srv/<app-name>/` and system configuration under `/etc/<app-name>/`. This makes deployments predictable for operators, compatible with distro tooling (backups, SELinux, package hooks), and easier to audit and automate. Use FHS-friendly packaging (no hidden state outside `/srv` and `/etc`) so services and maintenance scripts behave consistently.

## Web server configuration

- Document root: point the web server to `/srv/<app-name>` so each app has a single, self-contained document root.
- Config: store runtime config in `/etc/<app-name>/` and ensure services read that path.
