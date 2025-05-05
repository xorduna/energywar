# Energy War Game

## Configuration

The game supports multiple database configurations via environment variables:

### Database URI Options

- `DB_URI=game.db` (default SQLite file)
- `DB_URI=:memory:` (in-memory SQLite database)
- `DB_URI=file::memory:?cache=shared` (shared in-memory SQLite database)
- `DB_URI=file:./sqlite.db` (custom SQLite file path)
- `DB_URI=postgres://user:password@localhost:5432/dbname?sslmode=disable` (PostgreSQL connection)

### Server Port

- `SERVER_PORT=:8080` (default port)

### Example Usage

```bash
# Use default SQLite database
./energywar

# Use in-memory SQLite database
DB_URI=:memory: ./energywar

# Use PostgreSQL database
DB_URI=postgres://user:password@localhost:5432/energywar?sslmode=disable ./energywar

# Change server port
SERVER_PORT=:9000 ./energywar
```

## Features

- Flexible database configuration
- Environment variable-based settings
- Support for SQLite and PostgreSQL
- Easy deployment and configuration
