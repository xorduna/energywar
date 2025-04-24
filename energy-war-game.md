
**Description:**

This document describes a "battleship" like game but with power plants. Each user has to setup its power plants (nuclear, gas, wind or solar) to meet the required capacity. Each other players strike for turns to reduce its opponent capacity. If a plant is hit, all the capacity is removed.

Board 

- Letters for Y axis, and Number for X axis.

| Power plant | Code | Capacity | size  |
| ----------- | ---- | -------- | ----- |
| NUCLEAR     | N    | 1000     | 3 x 3 |
| GAS         | G    | 300      | 2 x 2 |
| WIND        | W    | 100      | 2 x 1 |
| SOLAR       | S    | 25       | 1 x 1 |
Other

| Code | Description |
| ---- | ----------- |
| -    | Land        |
| H    | Hit         |
| M    | Miss        |
| ?    | Unknown     |
|      |             |

Mechanics:
- User should build an energy infrastructure that meets at least the capacity defined in the game and max a 10% extra of the capacity
- If a power plant is HIT, capacity of the entire plant is removed from the counter
- The game ends when one of the players have below the 10% of the defined capacity

Example Board

```json
{
  "plants": [
    {
      "type": "NUCLEAR",
      "coordinates": ["A1", "A2", "A3", "B1", "B2", "B3"]
    },
    {
      "type": "GAS",
      "coordinates": ["C1", "C2", "D1", "D2"]
    },
    {
      "type": "WIND",
      "coordinates": ["E1", "E2"]
    },
    {
      "type": "SOLAR",
      "coordinates": ["F1"]
    },
  ],
  "hits": ["F1"],
  "misses": ["G1", "G2"],
  "total_capacity": 1425,
  "capacity": 1400
}
```


**Status mechanics:**

`PENDING` -> `IN_PROGRESS` -> `END`
(turn is handled by the `turn` field)

When a game starts, each player should post their initial board. They can update the board as many times as they want and once they are ready to play, the call `POST /games/:id/players/:name/ready`

Once both player have declared themselves ready, the game starts automatically until one of the two player wins.

**POST /games?size=10&capacity=1000&players=alice,bob**

- Starts a game
- Initial version will only support 2 players
- Initial status is pending
- Size should be between 5 and 20
- There is no limit for capacity but should be > 0
- turn is for the first player in alphabetical order

Response
```json
{
"id": "<id>",
"status": "PENDING",
"turn": "alice",
"winner": null,
"players": {
	"alice": {"ready": false, "total_capacity": 0, "capacity": 0},
	"bob": {"ready": false, "total_capacity": 0, "capacity": 0}
	}
}
```

Error Response
```json
{
"status": "ERROR",
"error": "INVALID_PARAMETERS"
}
```


**GET /games/:id**
- Gets the status of the game
Response example 1
```json
{
"id": "<id>",
"status": "PENDING",
"winner": null,
"players": {
	"alice": {"ready": false, "total_capacity": 0, "capacity": 0},
	"bob": {"ready": true, "total_capacity": 1000, "capacity": 1000}
	}
}
```

Response example 2
```json
{
"id": "<id>",
"status": "IN_PROGRESS",
"turn": "bob",
"winner": null,
"players": {
	"alice": {"ready": true, "total_capacity": 1000, "capacity": 700},
	"bob": {"ready": true, "total_capacity": 1000, "capacity": 800}
	}
}
```

Response example 3
```json
{
"id": "<id>",
"status": "END",
"turn": "alice",
"winner": "alice",
"players": {
	"alice": {"ready": true, "total_capacity": 1000, "capacity": 500},
	"bob": {"ready": true, "total_capacity": 1000, "capacity": 0}
	}
}
```

Status descriptions

| `game.status` | `winner`  | Description         |
| ------------- | --------- | ------------------- |
| `PENDING`     | `null`    | Incomplete setup    |
| `IN_PROGRESS` | `null`    | Game is in progress |
| `END`         | `"alice"` | Alice won           |
| `END`         | `"bob"`   | Bob won             |
The initial value of `turn` is the first alphabetically ordered player name

At the moment the game ends, the `turn` field still points to the player who performed the winning move (i.e., the winner).

**POST /games/:id/players/:name/ready**
- Sets the board set to true for this player
- To call this endpoint the game should be in `PENDING` and the capacity of the board should meet the expectations
- Once all players are ready, game starts and turn is for the first player in alphabetical order
- Clients should poll the endpoint /games/:id to know if they can strike

Response
```json
{"result": "OK"}
```

Error Response
```json
{
"status": "ERROR",
"error": "<error message>"
}
```

Possible errors:
- `GAME_ALREADY_STARTED`
- `BOARD_NOT_SET`


**POST /games/:id/players/:name/strike?target=target&y=A&x=1**
- player with name strikes player target
- After strike, game advances and turn is for next player
- Users can get the board to know the overal situation
- Check if name and target exist, otherwise return `INVALID_PLAYER`

Response
```json
{
"status": "OK",
"result": "HIT"
}
```

Possible values or result:
- `HIT`
- `MISS`

Error Response
```json
{
"status": "ERROR",
"error": "NOT_YOUR_TURN"
}
```

Possible error values
- `NOT_YOUR_TURN`
- `INVALID_COORDINATES`
- `INVALID_PLAYER`


**POST /games/:id/players/:name/board**
- Sets up board for player
- Two power plants cannot occupy the same coordinates
- Check rule capacities
- This endpoint should only be called while the status game is `PENDING`

Body:
```json
{
  "plants": [
    {
      "type": "NUCLEAR",
      "coordinates": ["A1", "A2", "A3", "B1", "B2", "B3"]
    },
    {
      "type": "GAS",
      "coordinates": ["C1", "C2", "D1", "D2"]
    },
    {
      "type": "WIND",
      "coordinates": ["E1", "E2"]
    },
    {
      "type": "SOLAR",
      "coordinates": ["F1"]
    },
  ],
}
```

Response
```json
{
  "plants": [
    {
      "type": "NUCLEAR",
      "coordinates": ["A1", "A2", "A3", "B1", "B2", "B3"]
    },
    {
      "type": "GAS",
      "coordinates": ["C1", "C2", "D1", "D2"]
    },
    {
      "type": "WIND",
      "coordinates": ["E1", "E2"]
    },
    {
      "type": "SOLAR",
      "coordinates": ["F1"]
    },
  ],
  "hits": ["F1"],
  "misses": ["G1", "G2"],
  "total_capacity": 1425,
  "capacity": 1400
}
```

Error response
```json
{
"status": "ERROR",
"error": "<error message>"
}
```

Possible errors:
- `OVERLAPPING_PLANTS`
- `INVALID_COORDINATES`
- `GAME_ALREADY_STARTED`

**GET /games/:id/players/:name/board**
- Get current user board
- Will require authentication eventually so it can only be requested by the user

Response
```json
{
  "plants": [
    {
      "type": "NUCLEAR",
      "coordinates": ["A1", "A2", "A3", "B1", "B2", "B3"]
    },
    {
      "type": "GAS",
      "coordinates": ["C1", "C2", "D1", "D2"]
    },
    {
      "type": "WIND",
      "coordinates": ["E1", "E2"]
    },
    {
      "type": "SOLAR",
      "coordinates": ["F1"]
    },
  ],
  "hits": ["F1"],
  "misses": ["G1", "G2"],
  "total_capacity": 1425,
  "capacity": 1400
}
```


**GET /games/:id/opponent/:name/board/map**
- Gets an ASCII representation of a players map
Response
```
   1 2 3 4 5 6 7 8 9 10
A  N N N . . . . . . .
B  N N N . . . . . . .
C  N N N . . . . . . .
D  . . . . . . . . . .
E  . . . . . . . . . .
F  . . . . . . . . . .
G  . . . . . . . . . .
H  . . . . . . . . . .
I  . . . . . . . . . .
J  . . . . . . . . . .

```

**GET /games/:id/opponent/:name/board**
- Gets a blind board (only hits and misses) for a player
- This endpoint is public
Response
```json
{
  "hits": ["B2"],
  "misses": ["A1", "C1"],
  "total_capacity": 1000,
  "capacity": 900,
}
```

**GET /games/:id/opponent/:name/board/map**
- Gets an ASCII representation of a blind map (only hits and misses for a player)
- This endpoint is public
Response
```
   1 2 3 4 5 6 7 8 9 10
A  M . . . . . . . . .
B  . H . . . . . . . .
C  M . . . . . . . . .
D  . . . . . . . . . .
E  . . . . . . . . . .
F  . . . . . . . . . .
G  . . . . . . . . . .
H  . . . . . . . . . .
I  . . . . . . . . . .
J  . . . . . . . . . .
