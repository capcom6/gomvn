<div align="center">
  <a href="https://github.com/capcom6/gomvn">
    <img src="assets/logo.png" alt="GoMVN logo" width="96" height="96">
  </a>

  <h1>GoMVN</h1>

  <p>A self-hosted HTTP repository for private Maven-format artifacts.</p>

  <p>
    <a href="https://capcom6.github.io/gomvn/">API reference</a>
    ·
    <a href="https://github.com/capcom6/gomvn/releases">Releases</a>
    ·
    <a href="https://github.com/capcom6/gomvn/issues">Issues</a>
  </p>
</div>

[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![License][license-shield]][license-url]

## Table of Contents

- [Table of Contents](#table-of-contents)
- [About the Project](#about-the-project)
- [Features](#features)
- [Built With](#built-with)
- [Getting Started](#getting-started)
  - [Container Image](#container-image)
  - [Release Archives](#release-archives)
- [Configuration](#configuration)
  - [Selecting a Configuration File](#selecting-a-configuration-file)
  - [Server and Permissions](#server-and-permissions)
  - [Database](#database)
  - [Repositories](#repositories)
  - [Storage](#storage)
- [First Start and User Management](#first-start-and-user-management)
- [Usage](#usage)
  - [Publishing with Gradle](#publishing-with-gradle)
  - [Consuming with Gradle](#consuming-with-gradle)
  - [Android and Other Clients](#android-and-other-clients)
- [Security](#security)
- [Development](#development)
- [Contributing](#contributing)
- [License](#license)
- [Contact](#contact)

## About the Project

GoMVN is a repository server for teams and organizations that need to publish and read private Java, Android, or other Maven-format artifacts on infrastructure they control. It provides artifact upload and download over HTTP, browsable repository indexes, token-based user access, and pluggable database and artifact storage backends.

The `release` and `snapshot` names used in the sample configuration are configurable route prefixes. GoMVN does not enforce release immutability or Maven snapshot-version policies.

## Features

- Upload artifacts with authenticated HTTP `PUT` requests
- Download artifacts and generated directory indexes with `GET`
- Serve configurable repository paths such as `release` and `snapshot`
- Manage repository users through a browser admin interface and API
- Generate user tokens and store only their bcrypt hashes
- Restrict each user to allowed repository paths and deploy permissions
- Use SQLite, MySQL, or PostgreSQL for application data
- Store artifacts on the local filesystem or in an S3-compatible bucket
- Use the standard AWS credential chain when S3 credentials are omitted
- Build multi-platform release archives and container images with GoReleaser

## Built With

- [Go](https://go.dev/) 1.26
- [Fiber](https://gofiber.io/)
- [GORM](https://gorm.io/)
- [Docker](https://www.docker.com/)
- [GoReleaser](https://goreleaser.com/)

## Getting Started

### Container Image

Prerequisites:

- Docker
- A reachable MySQL or PostgreSQL database
- A configuration file and a writable artifact directory

Copy the sample configuration and select an external database:

```bash
cp configs/config.example.yml config.yml
```

```yaml
database:
  driver: postgres
  dsn: host=postgres.example.com user=gomvn password=change-me dbname=gomvn port=5432 sslmode=require
```

Keep the remaining sample settings, including local artifact storage at `data/repository`. The database host must be reachable from the container; `localhost` inside the container refers to the container itself.

Create the data directory and start a pinned release image:

```bash
mkdir -p data

docker run -d \
  --name gomvn \
  --restart unless-stopped \
  --user "$(id -u):$(id -g)" \
  --publish 127.0.0.1:8080:8080 \
  --mount type=bind,source="$(pwd)/config.yml",target=/app/config.yml,readonly \
  --mount type=bind,source="$(pwd)/data",target=/app/data \
  ghcr.io/capcom6/gomvn:latest
```

Running the container as the host user makes the bind-mounted data directory writable. If S3 storage is configured instead, the directory is still created during database initialization, so keeping the mount avoids permission problems.

Read the generated administrator token from the logs:

```bash
docker logs gomvn
```

Then open `http://127.0.0.1:8080/admin/` or check the public index:

```bash
curl -fsS http://127.0.0.1:8080/
```

Pin a version tag or digest for repeatable deployments instead of relying on `latest`.

### Release Archives

GoReleaser publishes `tar.gz` archives for Linux and macOS and `zip` archives for Windows. Download the archive matching your operating system and CPU architecture from the [releases page](https://github.com/capcom6/gomvn/releases).

The current archive configuration packages the binary but does not declare the runtime `views` directory as an extra archive file. From a source checkout of the same release tag, place the matching configuration and views beside the extracted binary:

```bash
sudo mkdir -p /opt/gomvn
sudo chown "$(id -u):$(id -g)" /opt/gomvn
tar -xzf gomvn_Linux_x86_64.tar.gz -C /opt/gomvn

cp -R /path/to/gomvn-source/views /opt/gomvn/views
cp /path/to/gomvn-source/configs/config.example.yml /opt/gomvn/config.yml
```

Edit `/opt/gomvn/config.yml` to use MySQL or PostgreSQL, then start the server from the directory containing `views`:

```bash
cd /opt/gomvn
./gomvn --config /opt/gomvn/config.yml
```

On Windows, extract the ZIP and run `gomvn.exe` from the directory containing `config.yml` and `views`.

## Configuration

The application reads its main settings from YAML. Relative database and storage paths resolve from the process working directory.

### Selecting a Configuration File

The server uses this order:

1. The `--config` command-line flag
2. The `CONFIG_PATH` environment variable
3. `config.yml` in the process working directory

Examples:

```bash
CONFIG_PATH=/opt/gomvn/config.yml ./gomvn
./gomvn --config /opt/gomvn/config.yml
```

### Server and Permissions

| YAML path            | Type    | Sample or behavior                                      |
| -------------------- | ------- | ------------------------------------------------------- |
| `name`               | String  | Repository name shown by the index                      |
| `debug`              | Boolean | Enables GORM SQL logging; defaults to `false`           |
| `permissions.index`  | Boolean | Anonymous index access; defaults to `true` when omitted |
| `permissions.view`   | Boolean | Anonymous artifact reads; defaults to `false`           |
| `permissions.deploy` | Boolean | Anonymous artifact uploads; defaults to `false`         |
| `server.host`        | String  | Listen address; sample value is `0.0.0.0`               |
| `server.port`        | Integer | Listen port; sample value is `8080`                     |

Authenticated repository access is authorized against the user's allowed paths. The initial admin user has deploy access to `/`, so it can both publish and read all configured repository paths.

### Database

| YAML path         | Description                      |
| ----------------- | -------------------------------- |
| `database.driver` | `sqlite`, `mysql`, or `postgres` |
| `database.dsn`    | Driver-specific data source name |

### Repositories

`repository` is a list of URL path prefixes:

```yaml
repository:
  - release
  - snapshot
```

The names are user-defined. They do not automatically apply versioning or immutability rules.

### Storage

Choose `local` or `s3` with `storage.driver`.

For local storage:

```yaml
storage:
  driver: local
  options:
    root: data/repository
```

For S3-compatible storage:

```yaml
storage:
  driver: s3
  options:
    login: your-access-key
    password: your-secret-key
    region: us-east-1
    bucket: gomvn-artifacts
    prefix: repository
    endpoint: https://s3.example.com
```

Supported S3 options are:

| Option     | Description                                                                                                                                       |
| ---------- | ------------------------------------------------------------------------------------------------------------------------------------------------- |
| `login`    | Static access key; if either `login` or `password` is non-empty, the static credential provider takes precedence over AWS environment credentials |
| `password` | Static secret key used by the same static credential provider                                                                                     |
| `region`   | AWS region; overrides the region from the default AWS configuration when set                                                                      |
| `bucket`   | Bucket containing the repository tree                                                                                                             |
| `prefix`   | Key prefix inside the bucket                                                                                                                      |
| `endpoint` | Optional custom S3 endpoint; enables path-style addressing                                                                                        |

The AWS default credential chain is used only when both YAML `login` and `password` are empty. Setting either value activates the static credential provider and bypasses AWS credential environment variables. A configured YAML `endpoint` takes precedence over AWS SDK endpoint environment settings. Store static S3 credentials through a protected configuration file or secret mount rather than committing them.

## First Start and User Management

When the user table is empty, GoMVN creates an `admin` user with deploy access to `/` and prints its ID, username, and generated token to the server log. The plaintext token is not printed again because only its bcrypt hash is stored.

Use the administrator credentials to:

- Open the browser admin interface at `/admin/`
- List, create, update, and delete repository users
- Refresh a user's token
- Replace a user's allowed paths and per-path deploy permissions

Example authenticated API request:

```bash
curl -fsS \
  --user "admin:${GOMVN_TOKEN}" \
  http://127.0.0.1:8080/api/users
```

Example user creation:

```bash
curl -fsS \
  --user "admin:${GOMVN_TOKEN}" \
  --header 'Content-Type: application/json' \
  --data '{
    "name": "ci",
    "admin": false,
    "deploy": true,
    "allowed": ["/release/com/example/"]
  }' \
  http://127.0.0.1:8080/api/users
```

The create response contains the generated token. Store it as a secret and rotate it if it is exposed. See the [API reference](https://capcom6.github.io/gomvn/) and [`requests.http`](requests.http) for the complete management API examples.

## Usage

GoMVN exposes repositories as ordinary HTTP paths. Maven clients can use the repository URLs directly once credentials and repository names are configured.

The examples below use HTTPS because a TLS-terminating proxy is recommended for any non-local deployment.

### Publishing with Gradle

This Groovy example uses the standard `SNAPSHOT` suffix to select a repository path and reads credentials from Gradle properties or environment variables:

```groovy
def repositoryUser = findProperty("gomvnUsername") ?: System.getenv("GOMVN_USERNAME")
def repositoryToken = findProperty("gomvnToken") ?: System.getenv("GOMVN_TOKEN")

publishing {
    repositories {
        maven {
            name = "gomvn"
            url = uri(
                project.version.toString().endsWith("SNAPSHOT")
                    ? "https://gomvn.example.com/snapshot"
                    : "https://gomvn.example.com/release"
            )
            credentials {
                username = repositoryUser
                password = repositoryToken
            }
        }
    }

    publications {
        maven(MavenPublication) {
            from components.java
            groupId = "com.example"
            artifactId = "library"
            version = project.version
        }
    }
}
```

Do not commit the token. Supply it through a protected Gradle property, environment variable, or CI secret store.

### Consuming with Gradle

```groovy
def repositoryUser = findProperty("gomvnUsername") ?: System.getenv("GOMVN_USERNAME")
def repositoryToken = findProperty("gomvnToken") ?: System.getenv("GOMVN_TOKEN")

repositories {
    mavenCentral()
    maven {
        name = "gomvn"
        url = uri("https://gomvn.example.com/release")
        credentials {
            username = repositoryUser
            password = repositoryToken
        }
    }
}

dependencies {
    implementation "com.example:library:1.0.0"
}
```

### Android and Other Clients

The HTTP repository interface is independent of the build tool. Android projects can publish the component exposed by their Android Gradle Plugin version and consume artifacts from the same repository URLs. For other Maven-compatible clients, use HTTP Basic authentication with the generated username and token.

## Security

- GoMVN serves plain HTTP and does not configure TLS certificates. Terminate HTTPS at a reverse proxy or load balancer.
- Protect `config.yml`; database DSNs and static S3 credentials are stored as plain values.
- Treat generated user tokens as secrets. Only token hashes are persisted.
- Review anonymous access settings before starting the service. `permissions.view: true` and `permissions.deploy: true` bypass repository authentication for those operations.
- Give deploy users only the repository paths they need and disable deploy access for read-only consumers.
- Back up both the configured database and artifact storage.
- Do not expose the management API or admin interface to untrusted networks without additional access controls.

## Development

Prerequisites for a local source build:

- Go 1.26 or newer
- A C compiler when using the default SQLite driver
- `golangci-lint` for formatting and linting
- Air if using the optional live-reload target

Common commands:

```bash
go mod download
make fmt
make lint
make test
make build
```

Run the current source tree with SQLite and CGO:

```bash
cp configs/config.example.yml config.yml
CGO_ENABLED=1 go run . --config ./config.yml
```

Generate a local release snapshot with:

```bash
make release
```

## Contributing

1. Open an issue to confirm the intended behavior when it is not already tracked.
2. Create a focused branch from `master`.
3. Keep changes scoped and preserve existing project conventions.
4. Run `make lint`, `make test`, and `go build ./...`.
5. Open a pull request describing the user-visible behavior and verification performed.

## License

The repository is distributed under the MIT License. See [LICENSE](LICENSE).

## Contact

- Project and issue tracker: [capcom6/gomvn](https://github.com/capcom6/gomvn)
- Management API reference: [capcom6.github.io/gomvn](https://capcom6.github.io/gomvn/)

[contributors-shield]: https://img.shields.io/github/contributors/capcom6/gomvn.svg?style=for-the-badge
[contributors-url]: https://github.com/capcom6/gomvn/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/capcom6/gomvn.svg?style=for-the-badge
[forks-url]: https://github.com/capcom6/gomvn/network/members
[stars-shield]: https://img.shields.io/github/stars/capcom6/gomvn.svg?style=for-the-badge
[stars-url]: https://github.com/capcom6/gomvn/stargazers
[issues-shield]: https://img.shields.io/github/issues/capcom6/gomvn.svg?style=for-the-badge
[issues-url]: https://github.com/capcom6/gomvn/issues
[license-shield]: https://img.shields.io/github/license/capcom6/gomvn.svg?style=for-the-badge
[license-url]: https://github.com/capcom6/gomvn/blob/master/LICENSE
