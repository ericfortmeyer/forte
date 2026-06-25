# Roadmap

`forte` aims for wide adoption through intentionality and adherence to the Filesystem Hierarchy Standard (FHS). This roadmap outlines features prioritized before the 1.0.0 release.

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

- **Dry Run** – Safety through reliability, not simulation.
- **Mapping Flexibility** – Allow operators to customize destination paths (deferred until 5+ active users)
- **Permissions Sync** – Re-apply file permissions/ownership on redeploy
- **Cleanup & Archival** – Automated cleanup of old deployments and version history
- **Transactional Deploys** – Atomic swaps and rollback support

---

## Why is dry run deferred?

The goal is safety through reliability, not simulation. 0.3.0 focuses on bulletproof deployment and clear feedback rather than preview modes. Once real users report specific pain points, feature requests will be weighted against the core philosophy.

---

## Current Status

**Version:** 0.3.0
**Support:** Static sites, CGI-based, and runtime dependent deployments (e.g., PHP)

- Java
- NodeJS
- PHP
- Python
- Ruby (Rake)

| Source | Destination |
| -------- | ------------- |
| `/tmp/<appname>` | `/srv/<appname>` |
| `/tmp/<appname>-config` | `/etc/<appname>` |
| `/tmp/<appname>-assets` | `/srv/assets/<appname>/` |

---

## How to Contribute

- Feedback is welcome
- Please submit a pull request. See [CONTRIBUTING.md](./CONTRIBUTING.md)
