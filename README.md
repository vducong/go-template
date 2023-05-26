# Go template

Template for Golang services

## Content

- [Quick start](#quick-start)
- [Project structure](#project-structure)
- [Tools](#tools)

## Quick start

Local development:

```sh
# Install dependencies
$ make install

# Run app with migrations
$ make run
```

## Project structure

```shell
├── assets
│   └── credentials
├── cmd
│   └── app
├── configs
├── constants
├── githooks
├── internal
│   ├── app
│   ├── controller
│   ├── middleware
│   ├── [any domain]
│   │   ├── dto
│   │   ├── model
│   │   ├── module
│   │   ├── repo
│   │   └── service
│   ├── router
│   └── server
├── migrations
├── pkg
│   ├── databases
│   ├── failure
│   ├── logger
│   └── pubsub
├── scripts
│   └── bash
└── ...
```

### `assets`

Static files to go along with the repository.

### `cmd/app/main.go`

Configuration and logger initialization. Then the main function "continues" in
`internal/app/app.go`.

### `configs`

The config structure is in the `configs.go`.
First, default values of non-sensitive configuration are pre-defined within `default.go`.

We choose the [Viper](https://github.com/spf13/viper) library for reading config from `config.yaml`.
The environment variables are loaded if their tag match, then overwrite the default value.

### `constants`

Immutable values during the application's execution

### `githooks`

Git hooks

### `internal/app`

There is always one _Run_ function in the `app.go` file, which "continues" the _main_ function.

This is where all the main objects are created.
Dependency injection occurs through the "New ..." constructors.
This makes the business logic independent from other layers.

### `internal/controller`

Server handler layer for REST HTTP server [Gin framework](https://github.com/gin-gonic/gin).

### `internal/middleware`

Inspects and filters HTTP requests entering the application.

### `internal/router`

Server routers are written in the same style:

- Handlers are grouped by area of application (by a common basis)
- For each group, its own router structure is created, the methods of which process paths
- The structure of the business logic is injected into the router structure, which will be called by the handlers

### `internal/server`

Start the server and wait for signals for graceful completion.

### `internal/[domain]`

Demonstrating a sample of how to implement Domain Driven Design in Go.
Each package consists of:

- dto
- model: model of business logic can be used in any layer. There can also be internal methods of their own (e.g. validation).
- module: encapsulation of a domain's dependencies.
- repo: an abstract storage (database) that business logic works with. Under the hood, use [GORM](https://gorm.io/docs/index.html) for query and data manipulation.
- service: consists of business logic of the application.

### `migrations`

Contains SQL scripts used for database migration.

### `pkg/databases`

Initiate connections to all the kinds of databases.

### `pkg/failure`

Wrapper and transalator of application's errors.

### `pkg/logger`

We custom and utilize [zap](https://github.com/uber-go/zap) as our primary logger.

### `scripts`

Scripts to perform various build, install, analysis, etc. operations.

These scripts keep the root level Makefile small and simple.

## Tools

### Style checks

We use [golanglint-ci](https://golangci-lint.run/) to enforce a consistent code style across our codebase.
Place `.golanglint-ci.yml` at the root folder and follow the instruction to [integrate your editor](https://golangci-lint.run/usage/integrations/).
