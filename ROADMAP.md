# Roadmap

`forte` aims for wide adoption through intentionality and adherence to the Filesystem Hierarchy Standard (FHS). This roadmap outlines features prioritized before the 1.0.0 release.

---

## 0.3.0: Developer Experience

### 0.3.0 Features

- **Enhanced Documentation** – FHS rationale, web server templates (Apache, Caddy), framework-specific guides
- **Additional Examples** – Ruby (Rack), Java, and other common runtimes

### 0.3.0 Why This Matters

Operators need confidence in what `forte` will do before execution. Mutli-langauge example deployments and clearer docs reduce deployment anxiety.

#### Dry Run Deferred

Safety through reliability, not simulation. 0.3.0 focuses on bulletproof deployment and clear feedback rather than preview modes. Once real users report specific pain points, feature requests will be weighted against the core philosophy.

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

## Current Status

**Version:** 0.2.0
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
| `/tmp/<appname>-assets` *(0.2.0)* | `/srv/assets/<appname>/` *(0.2.0)* |

---

## How to Contribute

Feedback on the 0.2.0 design is welcome. See the [GitHub Epic]([#9](https://github.com/ericfortmeyer/forte/issues/9)) for discussion on static assets deployment.
