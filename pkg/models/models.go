package models

import (
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// GameStatus represents the status of a game
type GameStatus string

const (
	GameStatusPending    GameStatus = "PENDING"
	GameStatusInProgress GameStatus = "IN_PROGRESS"
	GameStatusEnd        GameStatus = "END"
)

// PlantType represents the type of a power plant
type PlantType string

const (
	PlantTypeNuclear PlantType = "NUCLEAR"
	PlantTypeGas     PlantType = "GAS"
	PlantTypeWind    PlantType = "WIND"
	PlantTypeSolar   PlantType = "SOLAR"
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

// PlayerInfo represents information about a player in a game
type PlayerInfo struct {
	Ready         bool   `json:"ready"`
	TotalCapacity int    `json:"total_capacity"`
	Capacity      int    `json:"capacity"`
	Token         string `json:"token"`
	Board         *Board `json:"board"`
}

// Game represents a game
type Game struct {
	gorm.Model
	ID       string                `json:"id" gorm:"primaryKey"`
	Status   GameStatus            `json:"status"`
	Turn     string                `json:"turn"`
	Winner   *string               `json:"winner,omitempty"`
	Players  map[string]PlayerInfo `json:"players" gorm:"serializer:json"`
	Size     int                   `json:"size"`
	Capacity int                   `json:"capacity"`
	Public   bool                  `json:"public"`
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

// PlantSize returns the size of a plant type
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

// ParseCoordinate parses a coordinate string into y and x values
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

// parseCoordinatePart parses a single part of a coordinate
func parseCoordinatePart(part string) (int, error) {
	if len(part) == 0 {
		return 0, errors.New("empty coordinate part")
	}

	if part[0] == 'A' {
		return int(part[1] - '0'), nil
	}

	return 0, fmt.Errorf("invalid coordinate part: %s", part)
}

// GenerateBlindBoard generates a blind board for an opponent
func (b *Board) GenerateBlindBoard() *Board {
	blindBoard := &Board{
		Plants:        make([]Plant, len(b.Plants)),
		Hits:          make([]string, len(b.Hits)),
		Misses:        make([]string, len(b.Misses)),
		TotalCapacity: b.TotalCapacity,
		Capacity:      b.Capacity,
	}

	copy(blindBoard.Hits, b.Hits)
	copy(blindBoard.Misses, b.Misses)

	for i, plant := range b.Plants {
		blindBoard.Plants[i] = Plant{
			Type:        plant.Type,
			Coordinates: make([]string, len(plant.Coordinates)),
		}
	}

	return blindBoard
}

// GenerateASCIIMap generates an ASCII representation of the board
func (b *Board) GenerateASCIIMap(size int, blind bool) string {
	grid := make([][]string, size)
	for i := range grid {
		grid[i] = make([]string, size)
		for j := range grid[i] {
			grid[i][j] = "."
		}
	}

	for _, plant := range b.Plants {
		for _, coord := range plant.Coordinates {
			y, x, _ := ParseCoordinate(coord)
			if blind {
				grid[y][x] = "?"
			} else {
				switch plant.Type {
				case PlantTypeNuclear:
					grid[y][x] = "N"
				case PlantTypeGas:
					grid[y][x] = "G"
				case PlantTypeWind:
					grid[y][x] = "W"
				case PlantTypeSolar:
					grid[y][x] = "S"
				}
			}
		}
	}

	for _, hit := range b.Hits {
		y, x, _ := ParseCoordinate(hit)
		grid[y][x] = "X"
	}

	for _, miss := range b.Misses {
		y, x, _ := ParseCoordinate(miss)
		grid[y][x] = "O"
	}

	var result strings.Builder
	for _, row := range grid {
		for _, cell := range row {
			result.WriteString(cell)
		}
		result.WriteString("\n")
	}

	return result.String()
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

// ErrorResponse represents an error response
type ErrorResponse struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}
