# 🏆 Trophies.gg

[![GitHub release (with filter)](https://img.shields.io/github/v/release/lukecarr/trophies)][release]
[![Docker image size](https://ghcr-badge.egpl.dev/lukecarr/trophies/size)][docker-images]
[![GitHub](https://img.shields.io/github/license/lukecarr/trophies)](LICENSE)
[![Code Climate maintainability](https://img.shields.io/codeclimate/maintainability/lukecarr/trophies)][codeclimate]
[![Matrix](https://img.shields.io/matrix/trophies%3Amatrix.org)][matrix]

> ⚠️ This project is currently pre-alpha and should be considered a WIP: breaking changes are to be expected frequently, and most functionality is missing/partially implemented!

Trophies.gg is a lightweight, self-hosted trophy tracker for PSN. The entire web app (including the backend, frontend, and persistent storage/database) is published as a single zero-dependency executable!

![image](https://github.com/lukecarr/trophies/assets/24438483/0e17e80f-c6bf-4c9a-bf6f-7082c938f90d)

## Installation

As mentioned, Trophies.gg is a single zero-dependency executable, so installation is as simple as downloading the latest binary from the project's [GitHub releases][releases] page for your system of choice.

### Docker

We also publish a Docker image that serves as a wrapper for the Trophies.gg binary.

```
docker run -d -p 3000:3000 ghcr.io/lukecarr/trophies:latest
```

You can find a comprehensive list of Docker image tags [here][docker-images].

### Examples

In the [examples](/examples) directory, you can find different scenarios for deploying/using Trophies.gg.

## Usage

The `trophies` binary is actually a CLI application that can perform many different operations (including running the all-in-one web server).

### `serve`

The `trophies serve` subcommand starts the Trophies.gg server. This is a combined server that includes both the backend and web-based frontend.

> When running this command, you might see a warning about "Launching in in-memory mode". Please check the README section on [Data persistence](#data-persistence) for further details.

### `migrate`

The `trophies migrate` subcommand performs SQL migrations. This command needs to be invoked after a fresh installation of Trophies.gg or an update that introduces new database requirements.

### `fetch`

This `trophies fetch` subcommand loads games, trophy lists, and trophy completion data from PSN into your local SQLite database.

## Data persistence

**By default, Trophies.gg launches in an "in-memory mode", where data does not persist across restarts.**

You should set the `DSN` environment variable to the path of an SQLite `.db` file (it will be automatically created if missing) so Trophies.gg can store persistent data.

### Warning message

To ensure that you're aware of the data persistence behaviour, a warning message will be logged on startup if you haven't set the `DSN` environment variable:

```txt
WARNING: Launching in in-memory mode as 'DSN' environment variable wasn't set. Data will be lost on shutdown!
```

## Environment variables

Trophies.gg uses environment variables to configure the application. Below, you can find details on the different environment variables that Trophies.gg looks for and what they are used to configure.

### `ADDR`

Indicates the address (and port) that Trophies.gg will listen on when calling `serve`.

The default address value is `:3000`.

### `DSN`

The connection string (or, more simply, path) of the SQLite database instance used to persist data.

If this variable isn't set, Trophies.gg will launch in "in-memory mode", and data will not persist across restarts.

### `DISABLE_IN_MEMORY_WARN`

Setting this variable to any value will disable the in-memory database warning message on startup.

### `NPSSO`

This is your PSN NPSSO token which is used to interact with the PSN API and fetch information on games, trophy lists, and trophy completion.

Your NPSSO token can be obtained by following these steps:

1. Visit the [PlayStation website][PlayStation] and sign in with your PSN account.
1. In the same browser session, access [this page][npsso] to obtain your NPSSO token.

The NPSSO token should return a response that looks like this:

```json
{ "npsso": "<64 character token>" }
```

Copy the token value (not including the quote characters) and use this as the value of the `NPSSO` environment variable.

### `RAWG_API_KEY`

The [RAWG API][rawg] is used by Trophies.gg to obtain accompanying game metadata, including background images, Metacritic scores, and genre information.

Using the link above, you can sign up for a free account and obtain an API key limited to 20K monthly requests.

> Trophies.gg caches results from the RAWG API in an effort to reduce API consumption (and keep you within their free tier limits).

Once you've obtained an API key, include it as an environment variable and enjoy beautiful game art and extra metadata!

## Screenshots

### Homepage (all games)

![image](https://github.com/lukecarr/trophies/assets/24438483/0e17e80f-c6bf-4c9a-bf6f-7082c938f90d)

### Individual game page

![image](https://github.com/lukecarr/trophies/assets/24438483/6e0444cc-4c87-4706-8809-fd7b2b2010a2)

## Versioning

Trophies.gg uses [Semantic Versioning v2.0.0][semver].

## License

Trophies.gg is distributed under the [Apache 2.0](LICENSE) license.

[release]: https://github.com/lukecarr/trophies/releases/latest
[codeclimate]: https://codeclimate.com/github/lukecarr/trophies
[matrix]: https://matrix.to/#/#trophies:matrix.org
[releases]: https://github.com/lukecarr/trophies/releases
[docker-images]: https://github.com/lukecarr/trophies/pkgs/container/trophies/versions
[PlayStation]: https://www.playstation.com/
[npsso]: https://ca.account.sony.com/api/v1/ssocookie
[semver]: https://semver.org/spec/v2.0.0.html
[rawg]: https://rawg.io/apidocs
