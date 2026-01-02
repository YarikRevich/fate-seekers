# fate-seekers

[![StandWithUkraine](https://raw.githubusercontent.com/vshymanskyy/StandWithUkraine/main/badges/StandWithUkraine.svg)](https://github.com/vshymanskyy/StandWithUkraine/blob/main/docs/README.md)

## General Information

The provided project is a demo used to demonstrate implemented opportunities of a game creation using Go native tools. Provided demo is an isometric world set in space.

## Features

There are implemented next logical components and systems:

* **Custom State Management**: A state management system based on ReduxJS approach.
* **Custom Renderer**: A rendering engine built from scratch.
* **Network Management**: Separation for latency-sensitive and non-latency sensitive API communication.
* **Input System**: Support for both gamepad and keyboard usage.
* **Sound Management**: Custom sound streams management for music and FX.
* **Collision System**: Custom collision detection implementation.
* **Interaction System**: Custom selectable items interaction system.
* **Monitoring System**: Custom monitoring system based on Grafana and Prometheus open-source tools.

## Setup

All setup related operations are processed via **Makefile** placed in the root directory.

### Build

In order to build the project, it's required to execute the next commands:

#### Client
```make
# Performs client build for operational version
make build-client-operational

# Performs client build for testing version
make build-client-testing
```

Different versions are used to have two independent clients on the same host(used for testing)

#### Server
```make
# Performs server build with included UI
make build-server-ui

# Performs server build with CLI only
make build-server-cli
```