# Siva - Fermion

Fermion is a user-friendly commerce platform that helps businesses manage their online and offline store (channel of sales), supply chains, finances, marketing, and other commercial operations through one streamlined dashboard. Fermion delivers a full suite of business management tools/apps. Essentials include product sourcing, sales and inventory tracking, payment processing, shipping, customer accounts, marketing, reporting, etc. Plus, you can expand your Fermion app stores/toolkit quickly with hundreds of Fermion Apps

# tech stack used

Golang
GORM
Go Python Interpreter
Event Driven Architecture (Kafka)

## Build & Run

```bash
# If u wanna run in development
$ ENV=DEV go run main.go

#if u wanna run in development with local db cleard
$ ENV=DEV go run main.go -clearDB=true

#if u wanna run in development with seed master data
$ ENV=DEV go run main.go -seedMD=true

#if u wanna run in development with migrations
$ ENV=DEV go run main.go -migrateDB=true

# If u wanna build
$ go build

#If u wanna run in docker
$ docker-compose up
```

- PORT : 3031

## Generate swagger documentation

install go swagger

```bash
go get github.com/swaggo/swag

cd $GOPATH/src/github.com/swaggo/swag

go install cmd/swag

```

to generate API documentation

```bash

swag init --parseDependency

```

to see the results, run app and access {{base_url}}/swagger/index.html

## [Manifest](manifest.md)

## [Ipaas](ipaas_core.md)
# Author

Siva-Fermion
