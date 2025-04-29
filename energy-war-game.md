# Energy War Game Documentation

## Game Concept

Energy War is a strategic multiplayer game where players compete by managing power plants and striking opponent's energy infrastructure.

## Core Mechanics

### Power Plants
- 4 Types of Power Plants:
  1. Nuclear (1000 MW, 3x3 grid)
  2. Gas (300 MW, 2x2 grid)
  3. Wind (100 MW, 2x1 grid)
  4. Solar (25 MW, 1x1 grid)

### Game Phases
1. **Game Creation**
   - Create game via index.html
   - Choose public or private game
   - Get shareable game link
   - Copy link to clipboard

2. **Game Joining**
   - Use join.html page
   - Enter game ID and player name
   - Receive unique authentication token

3. **Board Setup**
   - Players place power plants
   - Meet minimum capacity requirement (1000 MW)
   - Strategic plant placement
   - Validate board configuration

4. **Gameplay**
   - Turn-based strikes
   - Target opponent's power plants
   - Remove plant capacity when hit
   - Game ends when a player's capacity falls below 10%

## Authentication and Security

### Token-Based Authentication
- Unique token generated for each player
- Required for sensitive actions:
  * Setting board
  * Marking ready
  * Performing strikes
- Prevents unauthorized game manipulation

### Game Visibility
- Public games: Open to all
- Private games: Invite-only via shared link

## Frontend Components

### index.html
- Game creation interface
- Public/private game selection
- Shareable game link generation
- Clipboard copy functionality

### join.html
- Dedicated game joining page
- Supports URL-based joining
- Manual game ID and player name entry
- Input validation

### player.html
- Game board setup and interaction
- Real-time game state management
- Token-based authentication

## Technical Details

### Board Validation Rules
- Plants must fit within grid
- No overlapping plants
- Minimum total capacity requirement
- Maximum capacity limit

### Striking Mechanics
- One strike per turn
- Hits remove entire plant's capacity
- Misses end turn
- Game state updates in real-time

## Winning Conditions
- Reduce opponent's total capacity below 10%
- Strategic plant placement
- Efficient striking

## Future Enhancements
- Spectator mode
- Advanced matchmaking
- Persistent game storage
- Detailed game analytics

## Community and Feedback
- Open-source development
- Regular updates
- User-driven improvements
