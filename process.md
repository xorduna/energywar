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

Enhanced the frontend with interactive game features:
- Added a "New Game" button that creates a new game and displays the game ID
- Added a game viewing feature with an input field for game ID and a "View" button
- Added a play game section with input fields for game ID and player name
- Created a game.html page to display the game status and boards
- Implemented board visualization with power plant icons
- Added color coding for hits (red), misses (orange), and normal cells (green)
- Added visual indicators for working (green border) and damaged (red border) power plants
- Implemented automatic polling to update game status and boards every second
- Added error handling for API requests and null checks for board properties

Refactored the frontend for better maintainability:
- Organized frontend assets into a structured directory:
  ```
  frontend/assets/
  ├── css/        # Stylesheets
  ├── js/         # JavaScript files
  └── img/        # Image resources
  ```
- Separated concerns in the player.html page:
  - Moved JavaScript code to external file (player.js)
  - Moved CSS styles to external file (player.css)
  - Kept HTML structure clean and focused
- Improved code organization and maintainability
- Made the codebase more modular and easier to extend

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
- Multiplayer support

## Multiplayer Support

The game now supports multiple players with the following features:

1. **Player Management**
   - Players can join a game using a unique name
   - Each player has their own board and capacity

2. **Turn-Based Gameplay**
   - Players take turns in alphabetical order
   - Only the player whose turn it is can strike

3. **Opponent Boards**
   - Players can see blind versions of all opponent boards
   - Hits and misses are visible, but plant locations are hidden

4. **Frontend Enhancements**
   - The game view page shows boards for all players
   - The player page shows all opponent boards
   - Players can strike any opponent's board when it's their turn
