# Roadmap

`forte` aims for wide adoption through intentionality and adherence to the Filesystem Hierarchy Standard (FHS). This roadmap outlines features prioritized before the 1.0.0 release.

---

## 0.2.0: Core Deployment (Current)

### 0.2.0 Features

- **Static Assets Deployment** – Deploy static assets via `-assets` suffix (e.g., `/tmp/<appname>-assets` → `/srv/assets/<appname>/`)
- **Tarball Extraction** – Support `.tar.gz` archives as deployment sources
- **Runnable Docker Examples** – PHP (CGI), Python (WSGI/Flask), and Node.js (Express) examples with Nginx configuration demonstrating the FHS layout

### 0.2.0 Why This Matters

These features validate the core deployment model for web applications across multiple runtimes. Docker examples serve as end-to-end integration tests (via GitHub Actions healthchecks) and operational templates.

---

## 0.3.0: Developer Experience

### 0.3.0 Features

- **`--dry-run` Mode** – Preview deployment changes without applying them
- **Enhanced Documentation** – FHS rationale, web server templates (Apache, Caddy), framework-specific guides
- **Additional Examples** – Ruby (Rack), Java, and other common runtimes

### 0.3.0 Why This Matters

Operators need confidence in what `forte` will do before execution. Dry-run mode and clearer docs reduce deployment anxiety.

---

## 0.4.0: Distribution & Reach

### 0.4.0 Features

- **APT/Deb Packaging** – Install `forte` via standard Linux package managers
- **Binary Server Deployments** – Deploy Go binaries, compiled executables, and system services

### 0.4.0 Why This Matters

Package managers are the standard distribution mechanism for production tools. Binary support expands `forte` beyond web apps to system utilities.

---

## Post-1.0.0: Deferred Features

The following features are deferred until broader adoption and operator feedback:

- **Mapping Flexibility** – Allow operators to customize destination paths (deferred until 5+ active users)
- **Permissions Sync** – Re-apply file permissions/ownership on redeploy
- **Cleanup & Archival** – Automated cleanup of old deployments and version history
- **Transactional Deploys** – Atomic swaps and rollback support

---

## Current Status

**Version:** 0.1.0
**Support:** Static sites and CGI-based deployments (e.g., PHP)

| Source | Destination |
| -------- | ------------- |
| `/tmp/<appname>` | `/srv/<appname>` |
| `/tmp/<appname>-config` | `/etc/<appname>` |
| `/tmp/<appname>-assets` *(0.2.0)* | `/srv/assets/<appname>/` *(0.2.0)* |

---

## How to Contribute

Feedback on the 0.2.0 design is welcome. See the [GitHub Epic]([#9](https://github.com/ericfortmeyer/forte/issues/9)) for discussion on static assets deployment.
