# Forte Deployment Examples

Minimal deployment examples demonstrating how `forte deploy` installs services across different languages.

## Quick Start

Each example builds a Docker image and runs a minimal HTTP server. Build from the repository root:

```zsh
docker build -f examples/php/Dockerfile -t forte-example-php .
docker run --read-only --rm -d --name forte-example-php -p 8000:8000 forte-example-php
curl -fs --retry 5 --retry-delay 2 --retry-all-errors  http://localhost:8000/
docker stop forte-example-php
```
