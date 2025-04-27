# Energy War Game Architecture

## Overview

The Energy War Game is a multiplayer, turn-based strategy game implemented as a web application with a backend API and a frontend interface. The architecture is designed to be modular, scalable, and secure.

## System Components

### Backend (Go)
- **Language**: Go (Golang)
- **Web Framework**: Echo
- **API Documentation**: Swagger

#### Key Packages
1. `pkg/models`
   - Defines core data structures
   - Manages game and board representations
   - Handles coordinate and plant type validations

2. `pkg/game`
   - Implements game logic
   - Manages game state
   - Handles player actions and game progression

3. `pkg/handlers`
   - Manages API request handling
   - Implements request validation
   - Provides interface between API routes and game logic

### API Endpoints

#### Game Management
- `POST /games`: Create a new game
- `GET /games/:id`: Retrieve game status
- `GET /games/:id/status`: Get limited game information
- `POST /games/:id/join`: Join an existing game

#### Player Actions
- `POST /games/:id/players/:name/ready`: Mark player as ready
- `POST /games/:id/players/:name/board`: Set player's board
- `POST /games/:id/players/:name/strike`: Perform a strike action

#### Board Information
- `GET /games/:id/players/:name/board`: Get player's board
- `GET /games/:id/opponent/:name/board`: Get opponent's blind board

## Security Features

### Game Visibility
- Support for public and private game modes
- Configurable game visibility during game creation
- Limited information exposure for non-public games

### Token Management
- Player tokens are never exposed in game status endpoints
- Tokens only returned during game join process
- Consistent token hiding across all game information retrieval

## Game Mechanics

### Multiplayer Support
- 2-4 players per game
- Turn-based gameplay
- Alphabetical turn order
- Simultaneous board setup

### Power Plant Mechanics
- Four plant types: Nuclear, Gas, Wind, Solar
- Unique capacities and board sizes
- Strategic plant placement
- Capacity-based win conditions

## Frontend Architecture

### Technologies
- HTML5
- CSS3
- JavaScript (jQuery)
- Responsive design

### Key Features
- Interactive game creation
- Real-time game status updates
- Board visualization
- Player action management

## Security Principles

- No sensitive information exposure
- Consistent game state management
- Secure token handling
- Input validation at all levels

## Scalability Considerations

- Stateless API design
- Modular package structure
- Extensible game logic
- Easy to add new features or game modes

## Performance Optimization

- Efficient game state management
- Minimal data transfer
- Lightweight API responses
- Optimized game logic algorithms

## Deployment

- Containerization support
- Embedded frontend assets
- Single binary deployment
- Swagger documentation included

## Future Enhancements

- Persistent game storage
- Advanced matchmaking
- Spectator mode
- Detailed game analytics
