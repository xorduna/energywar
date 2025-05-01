package game

import (
	"errors"
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/xorduna/energywar/pkg/models"
	"gorm.io/gorm"
)

// GameManager manages all active games
type GameManager struct {
	games map[string]*models.Game
	mutex sync.RWMutex
	db    *gorm.DB
}

// NewGameManager creates a new game manager with a GORM connection
func NewGameManager(db *gorm.DB) *GameManager {
	return &GameManager{
		games: make(map[string]*models.Game),
		db:    db,
	}
}

// CreateGame creates a new game with the given parameters
func (gm *GameManager) CreateGame(size int, capacity int, public bool) (*models.Game, error) {
	// Validate parameters
	if size < 5 || size > 20 {
		return nil, errors.New("size should be between 5 and 20")
	}
	if capacity <= 0 {
		return nil, errors.New("capacity should be greater than 0")
	}

	// Create empty player info map
	playerInfoMap := make(map[string]models.PlayerInfo)

	// Generate a unique ID
	id := generateID()

	// Create the game
	game := &models.Game{
		ID:       id,
		Status:   models.GameStatusPending,
		Turn:     "",
		Winner:   nil,
		Players:  playerInfoMap,
		Size:     size,
		Capacity: capacity,
		Public:   public,
	}

	// Store the game in the database
	if err := gm.db.Create(game).Error; err != nil {
		return nil, err
	}

	// Store the game in memory
	gm.mutex.Lock()
	gm.games[id] = game
	gm.mutex.Unlock()

	return game, nil
}

// JoinGame allows a player to join an existing game
func (gm *GameManager) JoinGame(gameID string, playerName string) (string, error) {
	gm.mutex.Lock()
	defer gm.mutex.Unlock()

	// Get the game
	game, exists := gm.games[gameID]
	if !exists {
		return "", errors.New("game not found")
	}

	// Check if the game is still in PENDING status
	if game.Status != models.GameStatusPending {
		return "", errors.New("GAME_ALREADY_STARTED")
	}

	// Check if max players limit (4) is reached
	if len(game.Players) >= 4 {
		return "", errors.New("game is full (max 4 players)")
	}

	// Check if the player already exists
	if _, exists := game.Players[playerName]; exists {
		return "", errors.New("player already exists in this game")
	}

	// Generate a random token for the player
	token := generateToken()

	// Add the player to the game
	game.Players[playerName] = models.PlayerInfo{
		Ready:         false,
		TotalCapacity: 0,
		Capacity:      0,
		Token:         token,
		Board:         &models.Board{},
	}

	// If this is the first player, set the turn
	if game.Turn == "" && len(game.Players) > 0 {
		// Get all player names
		players := make([]string, 0, len(game.Players))
		for player := range game.Players {
			players = append(players, player)
		}
		// Sort players alphabetically
		sort.Strings(players)
		// Set turn to the first player
		game.Turn = players[0]
	}

	// Update the game in the database using a transaction
	err := gm.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(game).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return "", err
	}

	return token, nil
}

// GetGame retrieves a game by ID
func (gm *GameManager) GetGame(id string) (*models.Game, error) {
	gm.mutex.RLock()
	game, exists := gm.games[id]
	gm.mutex.RUnlock()

	if !exists {
		// Try to retrieve the game from the database using a transaction
		var dbGame models.Game
		err := gm.db.Transaction(func(tx *gorm.DB) error {
			if err := tx.First(&dbGame, "id = ?", id).Error; err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return nil, errors.New("game not found")
		}

		// Store the game in memory
		gm.mutex.Lock()
		gm.games[id] = &dbGame
		gm.mutex.Unlock()

		return &dbGame, nil
	}

	return game, nil
}

// SetBoard sets a player's board
func (gm *GameManager) SetBoard(gameID string, playerName string, board *models.Board) (*models.Board, error) {
	gm.mutex.Lock()
	defer gm.mutex.Unlock()

	// Get the game
	game, exists := gm.games[gameID]
	if !exists {
		return nil, errors.New("game not found")
	}

	// Check if the game is still in PENDING status
	if game.Status != models.GameStatusPending {
		return nil, errors.New("GAME_ALREADY_STARTED")
	}

	// Check if the player exists
	playerInfo, exists := game.Players[playerName]
	if !exists {
		return nil, errors.New("player not found")
	}

	// Validate the board
	if err := validateBoard(board, game.Size, game.Capacity); err != nil {
		return nil, err
	}

	// Calculate total capacity
	totalCapacity := 0
	for _, plant := range board.Plants {
		totalCapacity += models.PlantCapacity(plant.Type)
	}

	// Check if the total capacity meets the requirements
	if totalCapacity < game.Capacity || totalCapacity > int(float64(game.Capacity)*2) {
		return nil, fmt.Errorf("total capacity should be between %d and %d", game.Capacity, int(float64(game.Capacity)*1.1))
	}

	// Update the player's board
	board.TotalCapacity = totalCapacity
	board.Capacity = totalCapacity
	playerInfo.Board = board
	playerInfo.TotalCapacity = totalCapacity
	playerInfo.Capacity = totalCapacity
	game.Players[playerName] = playerInfo

	// Update the game in the database using a transaction
	err := gm.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(game).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return board, nil
}

// SetPlayerReady sets a player as ready
func (gm *GameManager) SetPlayerReady(gameID string, playerName string) error {
	gm.mutex.Lock()
	defer gm.mutex.Unlock()

	// Get the game
	game, exists := gm.games[gameID]
	if !exists {
		return errors.New("game not found")
	}

	// Check if the game is still in PENDING status
	if game.Status != models.GameStatusPending {
		return errors.New("GAME_ALREADY_STARTED")
	}

	// Check if the player exists
	playerInfo, exists := game.Players[playerName]
	if !exists {
		return errors.New("player not found")
	}

	// Check if the player has set their board
	if playerInfo.Board == nil || len(playerInfo.Board.Plants) == 0 {
		return errors.New("BOARD_NOT_SET")
	}

	// Set the player as ready
	playerInfo.Ready = true
	game.Players[playerName] = playerInfo

	// Check if all players are ready and there are at least 2 players
	allReady := true
	for _, info := range game.Players {
		if !info.Ready {
			allReady = false
			break
		}
	}

	// If all players are ready and there are at least 2 players, start the game
	if allReady && len(game.Players) >= 2 {
		game.Status = models.GameStatusInProgress

		// Set the turn to the first player in alphabetical order
		players := make([]string, 0, len(game.Players))
		for player := range game.Players {
			players = append(players, player)
		}
		sort.Strings(players)
		game.Turn = players[0]
	}

	// Update the game in the database using a transaction
	err := gm.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(game).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// Strike performs a strike action
func (gm *GameManager) Strike(gameID string, playerName string, targetName string, coord string) (string, error) {
	gm.mutex.Lock()
	defer gm.mutex.Unlock()

	// Get the game
	game, exists := gm.games[gameID]
	if !exists {
		return "", errors.New("game not found")
	}

	// Check if the game is in progress
	if game.Status != models.GameStatusInProgress {
		return "", errors.New("game is not in progress")
	}

	// Check if it's the player's turn
	if game.Turn != playerName {
		return "", errors.New("NOT_YOUR_TURN")
	}

	// Check if the player and target exist
	_, playerExists := game.Players[playerName]
	targetInfo, targetExists := game.Players[targetName]
	if !playerExists || !targetExists {
		return "", errors.New("INVALID_PLAYER")
	}

	// Validate the coordinate
	if err := models.ValidateCoordinate(coord, game.Size); err != nil {
		return "", errors.New("INVALID_COORDINATES")
	}

	// Check if the coordinate has already been hit or missed
	for _, hit := range targetInfo.Board.Hits {
		if hit == coord {
			return "", errors.New("coordinate already hit")
		}
	}
	for _, miss := range targetInfo.Board.Misses {
		if miss == coord {
			return "", errors.New("coordinate already missed")
		}
	}

	// Check if the coordinate hits a plant
	hit := false
	var hitPlantType models.PlantType
	var hitPlantCoords []string

	for _, plant := range targetInfo.Board.Plants {
		for _, plantCoord := range plant.Coordinates {
			if plantCoord == coord {
				hit = true
				hitPlantType = plant.Type
				hitPlantCoords = plant.Coordinates
				break
			}
		}
		if hit {
			break
		}
	}

	// Update the board based on the strike result
	if hit {
		// Add all plant coordinates to hits
		for _, plantCoord := range hitPlantCoords {
			if !contains(targetInfo.Board.Hits, plantCoord) {
				targetInfo.Board.Hits = append(targetInfo.Board.Hits, plantCoord)
			}
		}

		// Reduce capacity
		plantCapacity := models.PlantCapacity(hitPlantType)
		targetInfo.Capacity -= plantCapacity
		targetInfo.Board.Capacity -= plantCapacity

		// Check if the target has lost
		if targetInfo.Capacity <= int(float64(targetInfo.TotalCapacity)*0.1) {
			game.Status = models.GameStatusEnd
			winner := playerName
			game.Winner = &winner
		}
	} else {
		// Add to misses
		targetInfo.Board.Misses = append(targetInfo.Board.Misses, coord)
	}

	// Update the game state
	game.Players[targetName] = targetInfo

	// Update the turn if the game is still in progress
	if game.Status == models.GameStatusInProgress {
		// Find the next player
		players := make([]string, 0, len(game.Players))
		for player := range game.Players {
			players = append(players, player)
		}
		sort.Strings(players)

		// Find current player's index
		currentIndex := -1
		for i, player := range players {
			if player == playerName {
				currentIndex = i
				break
			}
		}

		// Set turn to next player (cycling through all players)
		nextIndex := (currentIndex + 1) % len(players)
		game.Turn = players[nextIndex]
	}

	// Update the game in the database using a transaction
	err := gm.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(game).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return "", err
	}

	if hit {
		return "HIT", nil
	}
	return "MISS", nil
}

// GetPlayerBoard retrieves a player's board
func (gm *GameManager) GetPlayerBoard(gameID string, playerName string) (*models.Board, error) {
	gm.mutex.RLock()
	game, exists := gm.games[gameID]
	gm.mutex.RUnlock()

	if !exists {
		// Try to retrieve the game from the database using a transaction
		var dbGame models.Game
		err := gm.db.Transaction(func(tx *gorm.DB) error {
			if err := tx.First(&dbGame, "id = ?", gameID).Error; err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return nil, errors.New("game not found")
		}

		// Store the game in memory
		gm.mutex.Lock()
		gm.games[gameID] = &dbGame
		gm.mutex.Unlock()

		game = &dbGame
	}

	// Check if the player exists
	playerInfo, exists := game.Players[playerName]
	if !exists {
		return nil, errors.New("player not found")
	}

	return playerInfo.Board, nil
}

// GetOpponentBlindBoard retrieves an opponent's blind board
func (gm *GameManager) GetOpponentBlindBoard(gameID string, opponentName string) (*models.Board, error) {
	gm.mutex.RLock()
	game, exists := gm.games[gameID]
	gm.mutex.RUnlock()

	if !exists {
		// Try to retrieve the game from the database using a transaction
		var dbGame models.Game
		err := gm.db.Transaction(func(tx *gorm.DB) error {
			if err := tx.First(&dbGame, "id = ?", gameID).Error; err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return nil, errors.New("game not found")
		}

		// Store the game in memory
		gm.mutex.Lock()
		gm.games[gameID] = &dbGame
		gm.mutex.Unlock()

		game = &dbGame
	}

	// Check if the opponent exists
	opponentInfo, exists := game.Players[opponentName]
	if !exists {
		return nil, errors.New("opponent not found")
	}

	// Generate a blind board
	return opponentInfo.Board.GenerateBlindBoard(), nil
}

// GetBoardMap generates an ASCII representation of a player's board
func (gm *GameManager) GetBoardMap(gameID string, playerName string, blind bool) (string, error) {
	gm.mutex.RLock()
	game, exists := gm.games[gameID]
	gm.mutex.RUnlock()

	if !exists {
		// Try to retrieve the game from the database using a transaction
		var dbGame models.Game
		err := gm.db.Transaction(func(tx *gorm.DB) error {
			if err := tx.First(&dbGame, "id = ?", gameID).Error; err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return "", errors.New("game not found")
		}

		// Store the game in memory
		gm.mutex.Lock()
		gm.games[gameID] = &dbGame
		gm.mutex.Unlock()

		game = &dbGame
	}

	// Check if the player exists
	playerInfo, exists := game.Players[playerName]
	if !exists {
		return "", errors.New("player not found")
	}

	// Generate the ASCII map
	return playerInfo.Board.GenerateASCIIMap(game.Size, blind), nil
}

// FormatGameStatus returns a string representation of the game status
func FormatGameStatus(game *models.Game) string {
	if game == nil {
		return "Game not found"
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("Game ID: %s\n", game.ID))
	result.WriteString(fmt.Sprintf("Status: %s\n", game.Status))

	if game.Status != models.GameStatusPending {
		result.WriteString(fmt.Sprintf("Turn: %s\n", game.Turn))
	}

	if game.Winner != nil {
		result.WriteString(fmt.Sprintf("Winner: %s\n", *game.Winner))
	}

	result.WriteString("Players:\n")
	for name, info := range game.Players {
		result.WriteString(fmt.Sprintf("- %s: Ready=%v, Capacity=%d/%d\n",
			name, info.Ready, info.Capacity, info.TotalCapacity))
	}

	return result.String()
}

// Helper functions

// generateID generates a unique ID for a game
func generateID() string {
	rand.Seed(time.Now().UnixNano())
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, 8)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// generateToken generates a random token for player authentication
func generateToken() string {
	rand.Seed(time.Now().UnixNano())
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 10)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// validateBoard validates a board configuration
func validateBoard(board *models.Board, size int, requiredCapacity int) error {
	// Check if the board has plants
	if len(board.Plants) == 0 {
		return errors.New("board has no plants")
	}

	// Create a 2D grid to check for overlapping plants
	grid := make([][]bool, size)
	for i := range grid {
		grid[i] = make([]bool, size)
	}

	// Check each plant
	for _, plant := range board.Plants {
		// Validate plant type
		switch plant.Type {
		case models.PlantTypeNuclear, models.PlantTypeGas, models.PlantTypeWind, models.PlantTypeSolar:
			// Valid plant type
		default:
			return fmt.Errorf("invalid plant type: %s", plant.Type)
		}

		// Get the expected size of the plant
		plantSize := models.PlantSize(plant.Type)
		expectedCoords := plantSize[0] * plantSize[1]

		// Check if the number of coordinates matches the expected size
		if len(plant.Coordinates) != expectedCoords {
			return fmt.Errorf("invalid number of coordinates for %s plant: expected %d, got %d", plant.Type, expectedCoords, len(plant.Coordinates))
		}

		// Check each coordinate
		for _, coord := range plant.Coordinates {
			fmt.Printf("Validating coordinate: %v\n", coord)
			// Validate the coordinate format
			if err := models.ValidateCoordinate(coord, size); err != nil {
				fmt.Printf("Error parsing coordinate %v: %v\n", coord, err)
				return fmt.Errorf("INVALID_COORDINATES: %v: %v", coord, err)
			}

			// Parse the coordinate
			y, x, err := models.ParseCoordinate(coord)
			if err != nil {
				return fmt.Errorf("INVALID_COORDINATES: %v: %v", coord, err)
			}

			// Check if the coordinate is within bounds
			if y < 0 || y >= size || x < 0 || x >= size {
				return errors.New("INVALID_COORDINATES: coordinate out of bounds")
			}

			// Check if the coordinate is already occupied
			if grid[y][x] {
				return errors.New("OVERLAPPING_PLANTS")
			}

			// Mark the coordinate as occupied
			grid[y][x] = true
		}

		// Validate plant shape
		if err := validatePlantShape(plant, size); err != nil {
			return err
		}
	}

	return nil
}

// validatePlantShape validates that a plant's coordinates form the correct shape
func validatePlantShape(plant models.Plant, size int) error {
	// Get the expected shape of the plant
	plantSize := models.PlantSize(plant.Type)
	width := plantSize[0]

	// Parse all coordinates
	coords := make([][2]int, len(plant.Coordinates))
	for i, coord := range plant.Coordinates {
		y, x, err := models.ParseCoordinate(coord)
		if err != nil {
			return err
		}
		coords[i] = [2]int{y, x}
	}

	// Sort coordinates by y, then by x
	sort.Slice(coords, func(i, j int) bool {
		if coords[i][0] == coords[j][0] {
			return coords[i][1] < coords[j][1]
		}
		return coords[i][0] < coords[j][0]
	})

	// Check if the coordinates form a rectangle of the correct size
	minY, minX := coords[0][0], coords[0][1]
	for i := 0; i < len(coords); i++ {
		expectedY := minY + (i / width)
		expectedX := minX + (i % width)

		if coords[i][0] != expectedY || coords[i][1] != expectedX {
			return fmt.Errorf("invalid plant shape for %s", plant.Type)
		}
	}

	return nil
}

// contains checks if a slice contains a string
func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}
