Snippet box App
======

Snippet box represents a standalone application that aims to provide necessary APIs to manage snippets domain.
The Snippet box domain provide necessary operation for user to view and add snippets.

How to launch locally
----

### Install necessary software

#### Software products and tools
* go >= 1.17
* PostgreSQL >= 11
* Docker Engine

#### Service's dependencies

* PostgresSQL

These services can be manually installed and configured or can be used semi-automatic scrips with [docker](https://docs.docker.com) engine.

**Note:** At the moment there are no instructions and/or script that provides support for manual installation. It is recommended to use docker engine.

### How to start locally

Prerequisites: all service dependencies should be up and running and properly configured.

* Build project:

```shell
go build -o main ./cmd/web/
```

* Show available flag:

```shell
./main -help
```

* Alternative - run on Docker:

```shell
docker-compose -f .\build\docker-compose.yaml up --build
```

### How to play

Exposed API can be checked using exposed [web](https://localhost:4000/) page

*Note* Http server is started on port 4000`

Design Consideration
----

TBA
