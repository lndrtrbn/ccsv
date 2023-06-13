# CCSV

Simple chat in GO (just to help me learn GO).

## Table of content

- [CCSV](#ccsv)
  - [Table of content](#table-of-content)
  - [Install dependencies](#install-dependencies)
  - [Build the application](#build-the-application)
  - [Flags](#flags)


## Install dependencies

```sh
go mod tidy
```

##  Build the application

```sh
# To build the server that receives and broadcasts the messages.
go build cmd/server/server.go
```

```sh
# To build the client UI that sends messages and display the conversation.
go build cmd/client/client.go
```

## Flags

For server:

| Name   | Usage                        | Default | Example       |
| ------ | ---------------------------- | ------- | ------------- |
| `port` | Port used by the HTTP server | `3000`  | `--port 3333` |

For client:

| Name     | Usage                           | Default                 | Example                          |
| -------- | ------------------------------- | ----------------------- | -------------------------------- |
| `name`   | Name of the user using the chat | none, mandatory         | `--name Michelle`                |
| `server` | Address of the HTTP server      | `http://localhost:3000` | `--server http://localhost:3333` |
