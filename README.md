<a id="readme-top"></a>

[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![project_license][license-shield]][license-url]

<!-- PROJECT LOGO -->
<br />
<div align="center">
  <a href="https://github.com/capcom6/gomvn">
    <img src="assets/logo.png" alt="Logo" width="80" height="80">
  </a>

  <h3 align="center">GoMVN</h3>

  <p align="center">
    A lightweight self-hosted repository manager for your private Maven artifacts.
    <br />
    <a href="https://github.com/capcom6/gomvn"><strong>Explore the docs »</strong></a>
    <br />
    <br />
    <a href="https://github.com/capcom6/gomvn/issues/new?labels=bug&template=bug-report---.md">Report Bug</a>
    ·
    <a href="https://github.com/capcom6/gomvn/issues/new?labels=enhancement&template=feature-request---.md">Request Feature</a>
  </p>
</div>

<!-- TABLE OF CONTENTS -->
- [About The Project](#about-the-project)
  - [Built With](#built-with)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
- [Configuration](#configuration)
  - [Storage](#storage)
- [User Guide](#user-guide)
- [Usage](#usage)
  - [How to create Java library](#how-to-create-java-library)
  - [How to create Android library](#how-to-create-android-library)
  - [How to use your private Maven repository](#how-to-use-your-private-maven-repository)
- [Roadmap](#roadmap)
- [Contributing](#contributing)
- [License](#license)
- [Contact](#contact)


<!-- ABOUT THE PROJECT -->
## About The Project

GoMVN is a lightweight, self-hosted Maven repository manager written in Go. It allows you to host your private Maven artifacts securely within your own infrastructure. With support for both release and snapshot repositories, user authentication, and flexible storage options, it's perfect for teams and organizations that need to manage their Java/Android libraries privately.

### Built With

* [![Go][Go.dev]][Go-url]
* [![Docker][Docker.com]][Docker-url]
* [![GORM][GORM.io]][GORM-url]

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- GETTING STARTED -->
## Getting Started

### Prerequisites

- Docker installed on your system
- A database (SQLite, MySQL, or PostgreSQL) - SQLite is used by default

### Installation

Use Docker to install this tool. The image is available at [GitHub](https://ghcr.io/capcom6/gomvn).

For better accessibility, map these Docker volumes:

| Path              | Description                                                                                  |
| ----------------- | -------------------------------------------------------------------------------------------- |
| `/app/data`       | app data for persistency                                                                     |
| `/app/config.yml` | configuration from outside of container, copy [default config](./configs/config.example.yml) |

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- CONFIGURATION -->
## Configuration

You can use `config.yml` to configure the service. Default `config.yml` is located in the same directory as the executable. To specify a different config file, use `--config /path/to/config.yml`.

Available configuration options:

| Path               | Description                                                            |
| ------------------ | ---------------------------------------------------------------------- |
| name               | name of the repository                                                 |
| debug              | enable debug output                                                    |
| permissions        | default permissions for the repository                                 |
| permissions.index  | anonymous access to index page                                         |
| permissions.view   | anonymous access to read artifacts                                     |
| permissions.deploy | anonymous access to deploy artifacts                                   |
| server             | http server configuration                                              |
| server.port        | port of the http server                                                |
| server.host        | host of the http server                                                |
| database           | database configuration                                                 |
| database.driver    | database driver to use (`sqlite`, `mysql`, `postgres`)                 |
| database.dsn       | database dsn, see https://gorm.io/docs/connecting_to_the_database.html |
| repository         | list of available repositories                                         |
| storage            | artifacts storage configuration                                        |
| storage.driver     | storage driver (`local` or `s3`)                                       |
| storage.options    | configuration options for storage driver, see below                    |

### Storage

Options for `local` storage driver:

- `root` - path to storage root.

Options for `s3` storage driver:

- `login` - user id/login;
- `password` - user secret/password;
- `endpoint`;
- `region`;
- `bucket`;
- `prefix`.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- USER GUIDE -->
## User Guide

On first run, admin account and his token is generated and printed into console.

You will need this to access [management API](https://capcom6.github.io/gomvn/) or local admin pages (http://my-private-repository.example.com/admin/), which is used to set user access.

If you don't have more users, you can use already created admin account to deploy and access your maven artifacts.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- USAGE -->
## Usage

### How to create Java library

Ensure that your `build.gradle` file contains configuration by this example:

```gradle
plugins {
    id 'maven-publish'
    id 'java'
}

publishing {
    repositories {
        maven {
            def releasesRepoUrl = "http://my-private-repository.example.com/release"
            def snapshotsRepoUrl = "http://my-private-repository.example.com/snapshot"
            name = 'mlj'
            url = project.version.endsWith('RELEASE') ? releasesRepoUrl : snapshotsRepoUrl
            credentials {
                username 'PUT HERE USERNAME'
                password 'PUT HERE TOKEN'
            }
        }
    }
    publications {
        maven(MavenPublication) {
            groupId = 'com.example'
            artifactId = 'library'
            version = '1.0.0.RELEASE'

            from components.java
        }
    }
}
```

### How to create Android library

Ensure that your `build.gradle` file contains configuration by this example:

```gradle
plugins {
    id 'maven-publish'
}

afterEvaluate {
    publishing {
        repositories {
            maven {
                def releasesRepoUrl = "http://my-private-repository.example.com/release"
                def snapshotsRepoUrl = "http://my-private-repository.example.com/snapshot"
                name = 'mlj'
                url = project.version.endsWith('RELEASE') ? releasesRepoUrl : snapshotsRepoUrl
                credentials {
                    username 'PUT HERE USERNAME'
                    password 'PUT HERE TOKEN'
                }
            }
        }
        publications {
            maven(MavenPublication) {
                // Applies the component for the release build variant.
                from components.release

                groupId = 'com.example'
                artifactId = 'library'
                version = '1.0.0.RELEASE'
            }
        }
    }
}
```

### How to use your private Maven repository

Append to your `build.gradle`:

```gradle
repositories {
    mavenCentral()
    maven {
        url "http://my-private-repository.example.com/release"
        credentials {
            username project.mljMavenUsername
            password project.mljMavenPassword
        }
    }
}

dependencies {
    implementation "com.example:library:1.0.0.RELEASE"
}
```

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- ROADMAP -->
## Roadmap

See the [open issues](https://github.com/capcom6/gomvn/issues) for a full list of proposed features (and known issues).

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- CONTRIBUTING -->
## Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement". Don't forget to give the project a star! Thanks again!

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- LICENSE -->
## License

Distributed under the MIT License. See `LICENSE` for more information.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- CONTACT -->
## Contact

Project Link: [https://github.com/capcom6/gomvn](https://github.com/capcom6/gomvn)

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- MARKDOWN LINKS & IMAGES -->
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
[Go.dev]: https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white
[Go-url]: https://go.dev/
[Docker.com]: https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white
[Docker-url]: https://www.docker.com/
[GORM.io]: https://img.shields.io/badge/GORM-003545?style=for-the-badge&logo=gorm&logoColor=white
[GORM-url]: https://gorm.io/
