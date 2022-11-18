GoMVN
=====

A lightweight self-hosted repository manager for your private Maven artifacts.


Installation
------------

Use docker to install this tool. Image is available at [Docker HUB](https://hub.docker.com/r/capcom6/gomvn).

For better accesibility, map these docker volumes:

| Path              | Description |
| ----------------- | ----------- |
| `/app/data`       | app data for persistency |
| `/app/config.yml` | configuration from outside of container, copy [default config](./configs/config.example.yml) |


### Configuration

You can use `config.yml` to configure the service. Default `config.yml` is located in the same directory as the executable. To specify a different config file, use `--config /path/to/config.yml`.

Available configuration options:

| Path               | Description |
| ------------------ | ----------- |
| name               | name of the repository |
| debug              | enable debug output |
| permissions        | default permissions for the repository |
| permissions.index  | anonymous access to index page |
| permissions.view   | anonymous access to read artifacts |
| permissions.deploy | anonymous access to deploy artifacts |
| server             | http server configuration |
| server.port        | port of the http server |
| server.host        | host of the http server |
| database           | database configuration |
| database.driver    | database driver to use (`sqlite`, `mysql`, `postgres`) |
| database.dsn       | database dsn, see https://gorm.io/docs/connecting_to_the_database.html |
| repository         | list of available repositories |


User Guide
----------

On first run, admin account and his token is generated and prited into console.

You will need this to access [management api](https://capcom6.github.io/gomvn/) or local admin pages (http://my-private-repository.example.com/admin/), which is used to set user access.

If you don't have more users, you can use already created admin accout to deploy and access your maven artifacts.


### How to create java library

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

#### How to create Android library

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

### How to use your private maven repository

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
