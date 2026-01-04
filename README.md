# fate-seekers

[![StandWithUkraine](https://raw.githubusercontent.com/vshymanskyy/StandWithUkraine/main/badges/StandWithUkraine.svg)](https://github.com/vshymanskyy/StandWithUkraine/blob/main/docs/README.md)

## General Information

The provided project is a **technical demo** used to demonstrate implemented opportunities of a game creation using Go native tools. Provided demo is an isometric world set in space.

## Features

The project implements the following logical components and systems:

* **State Management**: A centralized system to track the application's status, inspired by the Redux approach.
* **Rendering Pipeline**: A custom-built graphics engine designed to draw visuals from scratch.
* **Network Management**: Smart networking that separates real-time game data from standard API requests to reduce lag.
* **External Activity Interpolation**: A system that smooths out the movement of other players or objects to prevent visual jittering.
* **Debug IMGUI Menu**: An on-screen developer menu for testing, tweaking variables, and debugging while the app is running.
* **Multilayer Screen Transitions**: Handles smooth visual fading or sliding when moving between different menus or game screens.
* **Input System**: Full support for both keyboard/mouse and gamepad controllers.
* **Sound Management**: A custom audio engine that handles background music and sound effects separately.
* **Collision System**: A tailored system that detects when objects hit or overlap with each other.
* **Interaction System**: Logic that allows the user to select and interact with items in the world.
* **Events System**: An internal messaging system that lets different parts of the code talk to each other without being tightly connected.
* **Abstract Database Layering**: A flexible storage system that allows you to swap different database technologies without rewriting code.
* **Collectable Objects System**: Logic for spawning, picking up, and storing items in an inventory.
* **Global Notifications System**: A UI manager for displaying alerts, errors, or success messages to the user.
* **Subtitles System**: Handles the loading and display of text for dialogue or captions.
* **Runtime Settings**: Allows users to change video, audio, or gameplay settings instantly without restarting the application.
* **Translation System**: A localization framework for easily switching between different languages.
* **Monitoring System**: Performance tracking integrated with Grafana and Prometheus to keep an eye on system health.

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

### References

This technical demo utilizes a combination of public domain tilesets and AI-generated content.

#### Visual Assets
* Tileset: [Big Isometric Tileset (64x32)](https://opengameart.org/sites/default/files/big_isometric_tileset64x32.png) via OpenGameArt.
* AI Generation: Additional images and assets created using [PixelLabAI](https://www.pixellab.ai/).

#### Audio
* Sound Effects: Generated using [ElevenLabs](https://elevenlabs.io/).