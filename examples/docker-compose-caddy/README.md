# Docker Compose & Caddy

> This example demonstrates running Trophies.gg with [Docker Compose](https://docs.docker.com/compose/), using [Caddy](https://caddyserver.com/) to serve as a reverse proxy and TLS terminator.

## Requirements

- Docker 1.13.0+
- The `DOMAIN` environment variable set to the FQDN that you'd like to access Trophies.gg at.

## Getting Started

Download the [docker-compose.yml](docker-compose.yml) configuration file and pass it to `docker-compose`:

```bash
DOMAIN=example.com docker-compose -f /path/to/docker-compose.yml up -d
```

> The `-d` flag is passed to docker-compose so the containers are started in the background (detached).

🎉 **You now have Trophies.gg running on ports 80 and 443, serving HTTPS traffic from your configured domain!**
