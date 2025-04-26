# Energy War Game Architecture

## Overview

Energy War Game is a "battleship" like game where players set up power plants (nuclear, gas, wind, or solar) to meet a required capacity. Players take turns striking at each other's plants to reduce their opponent's capacity. The game ends when one of the players has below 10% of the defined capacity.

## System Architecture

The application follows a client-server architecture with a Go backend and HTML/CSS/JavaScript frontend.

### Backend

The backend is built using Go with the Echo framework and is organized into several packages:

```
energywar/
├── cmd/
│   └── server/
│       ├── main.go                # Server entry point
│       └── frontend/              # Frontend files embedded in the binary
├── pkg/
│   ├── models/                    # Data models
│   ├── game/                      # Game logic
│   └── handlers/                  # API handlers
└── docs/                          # Swagger documentation
```

#### Key Components

1. **Models (`pkg/models/models.go`)**
   - Defines data structures for the game (Plant, Board, Game, etc.)
   - Implements utility functions for coordinate parsing and validation
   - Provides board visualization functions

2. **Game Logic (`pkg/game/game.go`)**
   - Implements the `GameManager` to manage all active games
   - Handles game creation, board setup, player actions
   - Implements strike mechanics and win condition checking
   - Manages turn-based gameplay

3. **API Handlers (`pkg/handlers/handlers.go`)**
   - Implements RESTful API endpoints
   - Handles request validation and error responses
   - Provides Swagger annotations for API documentation

4. **Server (`cmd/server/main.go`)**
   - Sets up the Echo framework
   - Registers API routes
   - Configures middleware
   - Embeds frontend files
   - Sets up Swagger documentation

### Frontend

The frontend is built using HTML, CSS, and JavaScript with jQuery for AJAX requests. It follows a modular structure with separation of concerns:

```
frontend/
├── index.html              # Main landing page
├── game.html               # Game view page
├── player.html             # Player interaction page
└── assets/
    ├── css/                # Stylesheets
    │   └── player.css      # Styles for player page
    ├── js/                 # JavaScript files
    │   └── player.js       # Logic for player page
    └── img/                # Image resources
        ├── gas.png         # Gas plant icon
        ├── nuclear.png     # Nuclear plant icon
        ├── solar.png       # Solar plant icon
        └── wind.png        # Wind plant icon
```

It consists of three main pages:

1. **Index Page (`index.html`)**
   - Game description and rules
   - Power plant information
   - Game creation and joining interface

2. **Game View Page (`game.html`)**
   - Displays game status and player information
   - Shows boards for all players
   - Updates in real-time through polling

3. **Player Page (`player.html`)**
   - Board setup interface
   - Strike interface during gameplay
   - Real-time game status updates
   - Modular design with:
     - Separated CSS in player.css
     - Separated JavaScript in player.js

### API Endpoints

The API follows RESTful principles and is organized under the `/api` path:

#### Game Management
- `POST /api/games` - Create a new game
- `GET /api/games/:id` - Get game status
- `POST /api/games/:id/join` - Join a game

#### Player Actions
- `POST /api/games/:id/players/:name/ready` - Mark player as ready
- `POST /api/games/:id/players/:name/strike` - Strike a coordinate
- `POST /api/games/:id/players/:name/board` - Set player board
- `GET /api/games/:id/players/:name/board` - Get player board
- `GET /api/games/:id/players/:name/board/map` - Get ASCII representation of player board

#### Opponent Information
- `GET /api/games/:id/opponent/:name/board` - Get opponent blind board
- `GET /api/games/:id/opponent/:name/board/map` - Get ASCII representation of opponent blind board

## Data Flow

1. **Game Creation**
   - Client sends a request to create a game
   - Server creates a game with a unique ID
   - Client receives the game ID

2. **Player Joining**
   - Client sends a request to join a game with a player name
   - Server adds the player to the game
   - Client receives a token for authentication

3. **Board Setup**
   - Client sends a board configuration
   - Server validates the board
   - Client marks player as ready

4. **Gameplay**
   - Server determines turn order
   - Client sends strike requests
   - Server processes strikes and updates game state
   - Client polls for game updates

5. **Game End**
   - Server determines the winner
   - Client displays the result

## Multiplayer Support

The game supports multiple players with the following features:

1. **Player Management**
   - Players can join a game using a unique name
   - Each player has their own board and capacity

2. **Turn-Based Gameplay**
   - Players take turns in alphabetical order
   - Only the player whose turn it is can strike

3. **Opponent Boards**
   - Players can see blind versions of all opponent boards
   - Hits and misses are visible, but plant locations are hidden

4. **Real-time Updates**
   - Clients poll the server for game updates
   - UI updates to reflect the current game state

## Security Considerations

1. **Authentication**
   - Players receive a token when joining a game
   - Future implementation could require token for player-specific actions

2. **Input Validation**
   - All API inputs are validated
   - Coordinate validation ensures actions are within board boundaries
   - Plant placement validation prevents overlapping plants

## Future Enhancements

1. **WebSocket Support**
   - Replace polling with WebSocket for real-time updates

2. **User Authentication**
   - Implement user accounts and authentication

3. **Game History**
   - Store and display game history

4. **Enhanced UI**
   - Animations for strikes
   - Improved visual feedback

5. **Game Variations**
   - Different board sizes
   - Additional power plant types
   - Special abilities
