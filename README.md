# Energy War Game

A "battleship" like game but with power plants. Each user has to setup its power plants (nuclear, gas, wind or solar) to meet the required capacity. Players take turns striking to reduce their opponent's capacity. If a plant is hit, all the capacity is removed.

## Docker Deployment

### Building and Running Locally

1. Build the Docker image:
   ```bash
   docker build -t energywar .
   ```

2. Run the container:
   ```bash
   docker run -p 8080:8080 energywar
   ```

   Or using Docker Compose:
   ```bash
   docker-compose up
   ```

3. Access the application at http://localhost:8080


## Game Rules

| Power plant | Code | Capacity | size  |
| ----------- | ---- | -------- | ----- |
| NUCLEAR     | N    | 1000     | 3 x 3 |
| GAS         | G    | 300      | 2 x 2 |
| WIND        | W    | 100      | 2 x 1 |
| SOLAR       | S    | 25       | 1 x 1 |

### Mechanics
- User should build an energy infrastructure that meets at least the capacity defined in the game and max a 10% extra of the capacity
- If a power plant is HIT, capacity of the entire plant is removed from the counter
- The game ends when one of the players have below the 10% of the defined capacity

## API Documentation

The API documentation is available at `/swagger/index.html` when the application is running.

## Other documentation

- [architecture.md](architecture.md) - file containing other files description
- [energy-war-game.md](energy-war-game.md) - description for LLM development
- [process.md](process.md) - changes for LLM