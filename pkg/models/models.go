package models

import (
	"errors"
	"fmt"
	"strings"
)

// PlantType represents the type of power plant
type PlantType string

const (
	PlantTypeNuclear PlantType = "NUCLEAR"
	PlantTypeGas     PlantType = "GAS"
	PlantTypeWind    PlantType = "WIND"
	PlantTypeSolar   PlantType = "SOLAR"
)

// GameStatus represents the status of the game
type GameStatus string

const (
	GameStatusPending    GameStatus = "PENDING"
	GameStatusInProgress GameStatus = "IN_PROGRESS"
	GameStatusEnd        GameStatus = "END"
)

// Plant represents a power plant on the board
type Plant struct {
	Type        PlantType `json:"type"`
	Coordinates []string  `json:"coordinates"`
}

// Board represents a player's board
type Board struct {
	Plants        []Plant  `json:"plants"`
	Hits          []string `json:"hits"`
	Misses        []string `json:"misses"`
	TotalCapacity int      `json:"total_capacity"`
	Capacity      int      `json:"capacity"`
}

// PlayerInfo represents a player's information in the game
type PlayerInfo struct {
	Ready         bool   `json:"ready"`
	TotalCapacity int    `json:"total_capacity"`
	Capacity      int    `json:"capacity"`
	Token         string `json:"token,omitempty"`
	Board         *Board
}

// Game represents a game session
type Game struct {
	ID       string                `json:"id"`
	Status   GameStatus            `json:"status"`
	Turn     string                `json:"turn"`
	Winner   *string               `json:"winner"`
	Players  map[string]PlayerInfo `json:"players"`
	Size     int                   `json:"-"`
	Capacity int                   `json:"-"`
	Public   bool                  `json:"-"`
}

// PlantCapacity returns the capacity of a plant type
func PlantCapacity(plantType PlantType) int {
	switch plantType {
	case PlantTypeNuclear:
		return 1000
	case PlantTypeGas:
		return 300
	case PlantTypeWind:
		return 100
	case PlantTypeSolar:
		return 25
	default:
		return 0
	}
}

// PlantSize returns the size of a plant type as [width, height]
func PlantSize(plantType PlantType) [2]int {
	switch plantType {
	case PlantTypeNuclear:
		return [2]int{3, 3}
	case PlantTypeGas:
		return [2]int{2, 2}
	case PlantTypeWind:
		return [2]int{2, 1}
	case PlantTypeSolar:
		return [2]int{1, 1}
	default:
		return [2]int{0, 0}
	}
}

// ValidateCoordinate checks if a coordinate is valid for the given board size
func ValidateCoordinate(coord string, size int) error {
	if len(coord) < 2 {
		return errors.New("invalid coordinate format")
	}

	y := coord[0]
	x := coord[1:]

	// Check if y is a valid letter (A-Z)
	if y < 'A' || y > 'A'+byte(size-1) {
		return fmt.Errorf("y coordinate out of bounds: %c", y)
	}

	// Check if x is a valid number (1-size)
	var xVal int
	_, err := fmt.Sscanf(x, "%d", &xVal)
	if err != nil || xVal < 1 || xVal > size {
		return fmt.Errorf("x coordinate out of bounds: %s", x)
	}

	return nil
}

// ParseCoordinate converts a coordinate string (e.g., "A1") to [y, x] indices
func ParseCoordinate(coord string) (int, int, error) {
	if len(coord) < 2 {
		return 0, 0, errors.New("invalid coordinate format")
	}

	y := int(coord[0] - 'A')
	x := 0
	_, err := fmt.Sscanf(coord[1:], "%d", &x)
	if err != nil {
		return 0, 0, err
	}
	x-- // Convert to 0-based index

	return y, x, nil
}

// FormatCoordinate converts [y, x] indices to a coordinate string (e.g., "A1")
func FormatCoordinate(y, x int) string {
	return fmt.Sprintf("%c%d", 'A'+byte(y), x+1)
}

// GenerateBlindBoard generates a board with only hits and misses visible
func (b *Board) GenerateBlindBoard() *Board {
	return &Board{
		Hits:          b.Hits,
		Misses:        b.Misses,
		TotalCapacity: b.TotalCapacity,
		Capacity:      b.Capacity,
	}
}

// GenerateASCIIMap generates an ASCII representation of the board
func (b *Board) GenerateASCIIMap(size int, blind bool) string {
	// Create a 2D grid
	grid := make([][]string, size)
	for i := range grid {
		grid[i] = make([]string, size)
		for j := range grid[i] {
			grid[i][j] = "."
		}
	}

	// Place plants if not blind
	if !blind {
		for _, plant := range b.Plants {
			for _, coord := range plant.Coordinates {
				y, x, err := ParseCoordinate(coord)
				if err != nil {
					continue
				}
				if y >= 0 && y < size && x >= 0 && x < size {
					grid[y][x] = string(plant.Type[0])
				}
			}
		}
	}

	// Place hits
	for _, coord := range b.Hits {
		y, x, err := ParseCoordinate(coord)
		if err != nil {
			continue
		}
		if y >= 0 && y < size && x >= 0 && x < size {
			grid[y][x] = "H"
		}
	}

	// Place misses
	for _, coord := range b.Misses {
		y, x, err := ParseCoordinate(coord)
		if err != nil {
			continue
		}
		if y >= 0 && y < size && x >= 0 && x < size {
			grid[y][x] = "M"
		}
	}

	// Generate the ASCII map
	var sb strings.Builder

	// Header row with column numbers
	sb.WriteString("   ")
	for i := 1; i <= size; i++ {
		sb.WriteString(fmt.Sprintf("%d ", i))
	}
	sb.WriteString("\n")

	// Board rows
	for i := 0; i < size; i++ {
		sb.WriteString(fmt.Sprintf("%c  ", 'A'+byte(i)))
		for j := 0; j < size; j++ {
			sb.WriteString(grid[i][j] + " ")
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

// StrikeResponse represents a strike response
type StrikeResponse struct {
	Status string `json:"status"`
	Result string `json:"result"`
}

// ReadyResponse represents a ready response
type ReadyResponse struct {
	Result string `json:"result"`
}
