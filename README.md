# 🏆 Trophies.gg

> ⚠️ This project is currently pre-alpha and should be considered a WIP: breaking changes are to be expected frequently, and most functionality is missing/partially implemented!

Trophies.gg is a lightweight, self-hosted trophy tracker for PSN. The entire web app (including the backend, frontend, and persistent storage/database) is published as a single zero-dependency Docker image!

[![GitHub release (with filter)](https://img.shields.io/github/v/release/lukecarr/trophies)][release]
[![Docker image size](https://ghcr-badge.egpl.dev/lukecarr/trophies/size)][docker-images]
[![GitHub](https://img.shields.io/github/license/lukecarr/trophies)](LICENSE)
[![Code Climate maintainability](https://img.shields.io/codeclimate/maintainability/lukecarr/trophies)][codeclimate]
[![Matrix](https://img.shields.io/matrix/trophies%3Amatrix.org)][matrix]
[![Open in Gitpod](https://img.shields.io/badge/open_in-gitpod-FFAE33?logo=gitpod)][gitpod]

## Installation

Trophies.gg is distributed as a Docker image. The steps involve performing SQL migrations (this needs to happen on a fresh install or when upgrading to a new version), fetching data from PSN, and then starting the server:

```shell
# Pull the Trophies.gg Docker image
$ docker pull ghcr.io/lukecarr/trophies

# Setup an environment variable pointing to where your Trophies.gg database will persist
$ export DATABASE_PATH=/path/to/your/folder

# Perform migrations
$ docker run -it -v $DATABASE_PATH:/data ghcr.io/lukecarr/trophies migrate

# Obtain an NPSSO token (see README section) from PSN
$ export NPSSO=<your NPSSO token>

# Fetch data from PSN
$ docker run -it -e NPSSO=$NPSSO -v $DATABASE_PATH:/data ghcr.io/lukecarr/trophies fetch

# Obtain a RAWG API key (if you want game metadata, background images, etc.)
$ export RAWG_API_KEY=<your API key>

# Launch the server on port 3000 (-p 3000:3000) in detached mode (-d)
$ docker run -d -e NPSSO=$NPSSO -e RAWG_API_KEY=$RAWG_API_KEY -v $DATABASE_PATH:/data ghcr.io/lukecarr/trophies
```

> You can find a comprehensive list of Docker image tags [here][docker-images].

### Examples

In the [examples](/examples) directory, you can find different scenarios for deploying Trophies.gg.

## Usage

The `trophies` binary (the Docker image entrypoint) is actually a CLI application that can perform many different operations (including running the all-in-one web server).

### `serve`

The `trophies serve` subcommand starts the Trophies.gg server. This is a combined server that includes both the backend and web-based frontend.

**This subcommand is the default Dockerfile command.**

### `migrate`

The `trophies migrate` subcommand performs SQL migrations. This command needs to be invoked after a fresh installation of Trophies.gg or an update that introduces new database requirements.

### `fetch`

This `trophies fetch` subcommand loads games, trophy lists, and trophy completion data from PSN.

## Environment variables

Trophies.gg uses environment variables to configure the application. Below, you can find details on the different environment variables that Trophies.gg looks for and what they are used to configure.

### `NPSSO`

This is your PSN NPSSO token which is used to interact with the PSN API and fetch information on games, trophy lists, and trophy completion.

Your NPSSO token can be obtained by following these steps:

1. Visit the [PlayStation website][PlayStation] and sign in with your PSN account.
1. Cccess [this page][npsso] in the same browser session to obtain your NPSSO token.

The NPSSO token should return a response that looks like this:

```json
{ "npsso": "<64 character token>" }
```

Copy the token value (not including the quote characters) and use this as the value of the `NPSSO` environment variable.

### `RAWG_API_KEY`

The [RAWG API][rawg] is used by Trophies.gg to obtain accompanying game metadata, including background images, Metacritic scores, and genre information.

You can sign up for a free account using the link above and obtain an API key limited to 20K monthly requests.

> Trophies.gg caches results from the RAWG API to reduce API consumption (and keep you within their free tier limits).

Once you've obtained an API key, include it as an environment variable and enjoy beautiful game art and extra metadata!

### `DISABLE_NEW_VERSION_CHECK`

By default, Trophies.gg will make a request to GitHub's API on startup to check for new releases.

If you wish to disable this logic, please set the `DISABLE_NEW_VERSION_CHECK` variable to any value.

## Screenshots

### Homepage (all games)

![image](https://github.com/lukecarr/trophies/assets/24438483/20e5ae31-8d3c-45e6-8f6c-15e973811e8f)

### Individual game page

![image](https://github.com/lukecarr/trophies/assets/24438483/6e0444cc-4c87-4706-8809-fd7b2b2010a2)

## Versioning

Trophies.gg uses [Semantic Versioning v2.0.0][semver].

## License

Trophies.gg is distributed under the [Apache 2.0](LICENSE) license.

[release]: https://github.com/lukecarr/trophies/releases/latest
[codeclimate]: https://codeclimate.com/github/lukecarr/trophies
[matrix]: https://matrix.to/#/#trophies:matrix.org
[gitpod]: https//gitpod.io/#https://github.com/lukecarr/trophies
[releases]: https://github.com/lukecarr/trophies/releases
[docker-images]: https://github.com/lukecarr/trophies/pkgs/container/trophies/versions
[PlayStation]: https://www.playstation.com/
[npsso]: https://ca.account.sony.com/api/v1/ssocookie
[semver]: https://semver.org/spec/v2.0.0.html
[rawg]: https://rawg.io/apidocs
