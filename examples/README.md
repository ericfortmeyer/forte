# Forte Deployment Examples

Minimal deployment examples demonstrating how `forte deploy` installs services across different languages.

> ⚠️ EXAMPLE ONLY - NOT TO BE USED IN PRODUCTION
> The Dockerfiles in the subdirectories demonstrate deployments using <https://github.com/ericfortmeyer/forte>
> Production deployments required additional hardening: resource limits, server configuration, resource limits, secrets management, and proper logging.

## Quick Start

Each example builds a Docker image and runs a minimal HTTP server. Build from the repository root:

```zsh
docker build -f examples/php/Dockerfile -t forte-example-php .
docker run --read-only --rm -d --name forte-example-php -p 8000:8000 forte-example-php
curl -fs --retry 5 --retry-delay 2 --retry-all-errors  http://localhost:8000/
docker stop forte-example-php
```
