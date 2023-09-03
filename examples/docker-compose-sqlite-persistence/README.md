# Docker Compose with Persistence

> This example demonstrates running Trophies.gg with [Docker Compose](https://docs.docker.com/compose/), using a [Docker volume](https://docs.docker.com/storage/volumes/) to store the SQLite database (for data persistence).

## Requirements

* Docker 1.13.0+

## Getting Started

Download the [docker-compose.yml](docker-compose.yml) configuration file and pass it to `docker-compose`:

```bash
docker-compose -f /path/to/docker-compose.yml up -d
```

> The `-d` flag is passed to docker-compose so the container is started in the background (detached).

🎉 **You now have Trophies.gg running on port 3000, and your data will persist between restarts!**

> If you'd like to add automatic, zero-configuration HTTPS to your site, take a look at our [Docker Compose & Caddy example](../docker-compose-caddy/README.md)!
