# Energy War Game API Implementation Process

## Overview

This document describes the implementation process of the Energy War Game API, a "battleship" like game but with power plants. The implementation follows the requirements specified in the energy-war-game.md file.

## Implementation Steps

### 1. Project Setup

- Created a Go module for the project
```bash
go mod init github.com/xorduna/energywar
```

- Installed required dependencies
```bash
go get github.com/labstack/echo/v4
go get github.com/swaggo/echo-swagger
go get github.com/swaggo/swag/cmd/swag
```

- Created the project directory structure
```bash
mkdir -p cmd/server pkg/models pkg/handlers pkg/game frontend docs
```

### 2. Data Models Implementation

Created the data models in `pkg/models/models.go`:
- Defined power plant types (NUCLEAR, GAS, WIND, SOLAR)
- Defined game status types (PENDING, IN_PROGRESS, END)
- Implemented board and game structures
- Added utility functions for coordinate parsing and validation
- Implemented board visualization functions

### 3. Game Logic Implementation

Created the game logic in `pkg/game/game.go`:
- Implemented the `GameManager` to manage all active games
- Added functions for creating games, setting up boards, and handling player actions
- Implemented strike mechanics and win condition checking
- Added validation for board setup and plant placement
- Implemented turn management

### 4. API Handlers Implementation

Created the API handlers in `pkg/handlers/handlers.go`:
- Implemented all required endpoints as specified in the requirements
- Added Swagger annotations for API documentation
- Implemented request validation and error handling
- Added JSON response formatting

### 5. Frontend Implementation

Created a simple frontend in `frontend/index.html`:
- Added game description and rules
- Included power plant information
- Added API usage examples
- Provided a link to the Swagger documentation

### 6. Server Implementation

Created the server in `cmd/server/main.go`:
- Set up the Echo framework
- Registered all API routes
- Configured middleware (logging, CORS, etc.)
- Embedded frontend files
- Added Swagger documentation setup

### 7. Build System

Created a Makefile with the following targets:
- `build`: Builds the binary and generates Swagger documentation
- `run`: Builds and runs the application
- `swagger`: Generates Swagger documentation
- `clean`: Cleans build artifacts
- `test`: Runs tests
- `deps`: Installs dependencies

## API Endpoints

The following API endpoints were implemented:

- `POST /api/games`: Create a new game
- `GET /api/games/:id`: Get game status
- `POST /api/games/:id/players/:name/ready`: Set player as ready
- `POST /api/games/:id/players/:name/strike`: Strike a coordinate
- `POST /api/games/:id/players/:name/board`: Set player board
- `GET /api/games/:id/players/:name/board`: Get player board
- `GET /api/games/:id/players/:name/board/map`: Get ASCII representation of player board
- `GET /api/games/:id/opponent/:name/board`: Get opponent blind board
- `GET /api/games/:id/opponent/:name/board/map`: Get ASCII representation of opponent blind board

## Final Result

The Energy War Game API is running at http://localhost:8080 with API endpoints under the `/api` path. The frontend is served from the root path `/` and the Swagger documentation is available at http://localhost:8080/swagger/index.html.

The implementation successfully meets all the requirements:
- Backend in Go latest version
- Using the Echo framework
- Swagger documentation provided
- Makefile for building binary and documentation
- Frontend files embedded in the binary

The game mechanics are fully implemented according to the specifications, including:
- Power plant setup with different types, capacities, and sizes
- Turn-based gameplay
- Strike mechanics
- Win condition checking
