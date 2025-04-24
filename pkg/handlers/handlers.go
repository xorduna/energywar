package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/xorduna/energywar/pkg/game"
	"github.com/xorduna/energywar/pkg/models"
)

// Handler contains all the handler functions for the API
type Handler struct {
	GameManager *game.GameManager
}

// NewHandler creates a new handler
func NewHandler(gm *game.GameManager) *Handler {
	return &Handler{
		GameManager: gm,
	}
}

// @Summary Create a new game
// @Description Creates a new game with the specified parameters
// @Tags games
// @Accept json
// @Produce json
// @Param size query int false "Board size (5-20)" default(10)
// @Param capacity query int false "Required capacity" default(1000)
// @Param players query string true "Comma-separated list of player names"
// @Success 200 {object} models.Game
// @Failure 400 {object} models.ErrorResponse
// @Router /games [post]
func (h *Handler) CreateGame(c echo.Context) error {
	// Parse query parameters
	sizeStr := c.QueryParam("size")
	capacityStr := c.QueryParam("capacity")
	playersStr := c.QueryParam("players")

	// Default values
	size := 10
	capacity := 1000

	// Parse size
	if sizeStr != "" {
		var err error
		size, err = strconv.Atoi(sizeStr)
		if err != nil || size < 5 || size > 20 {
			return c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Status: "ERROR",
				Error:  "INVALID_PARAMETERS",
			})
		}
	}

	// Parse capacity
	if capacityStr != "" {
		var err error
		capacity, err = strconv.Atoi(capacityStr)
		if err != nil || capacity <= 0 {
			return c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Status: "ERROR",
				Error:  "INVALID_PARAMETERS",
			})
		}
	}

	// Parse players
	if playersStr == "" {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status: "ERROR",
			Error:  "INVALID_PARAMETERS",
		})
	}
	players := strings.Split(playersStr, ",")
	if len(players) != 2 {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status: "ERROR",
			Error:  "INVALID_PARAMETERS",
		})
	}

	// Create the game
	gameObj, err := h.GameManager.CreateGame(size, capacity, players)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status: "ERROR",
			Error:  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, gameObj)
}

// @Summary Get game status
// @Description Gets the current status of a game
// @Tags games
// @Accept json
// @Produce json
// @Param id path string true "Game ID"
// @Success 200 {object} models.Game
// @Failure 404 {object} models.ErrorResponse
// @Router /games/{id} [get]
func (h *Handler) GetGame(c echo.Context) error {
	// Get game ID from path
	id := c.Param("id")

	// Get the game
	gameObj, err := h.GameManager.GetGame(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, models.ErrorResponse{
			Status: "ERROR",
			Error:  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, gameObj)
}

// @Summary Set player ready
// @Description Sets a player as ready to start the game
// @Tags players
// @Accept json
// @Produce json
// @Param id path string true "Game ID"
// @Param name path string true "Player name"
// @Success 200 {object} models.ReadyResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /games/{id}/players/{name}/ready [post]
func (h *Handler) SetPlayerReady(c echo.Context) error {
	// Get game ID and player name from path
	id := c.Param("id")
	name := c.Param("name")

	// Set the player as ready
	err := h.GameManager.SetPlayerReady(id, name)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status: "ERROR",
			Error:  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.ReadyResponse{
		Result: "OK",
	})
}

// @Summary Strike a coordinate
// @Description Player strikes a coordinate on the opponent's board
// @Tags players
// @Accept json
// @Produce json
// @Param id path string true "Game ID"
// @Param name path string true "Player name"
// @Param target query string true "Target player name"
// @Param y query string true "Y coordinate (A-Z)"
// @Param x query int true "X coordinate (1-size)"
// @Success 200 {object} models.StrikeResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /games/{id}/players/{name}/strike [post]
func (h *Handler) Strike(c echo.Context) error {
	// Get game ID and player name from path
	id := c.Param("id")
	name := c.Param("name")

	// Get query parameters
	target := c.QueryParam("target")
	y := c.QueryParam("y")
	xStr := c.QueryParam("x")

	// Validate parameters
	if target == "" || y == "" || xStr == "" {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status: "ERROR",
			Error:  "INVALID_PARAMETERS",
		})
	}

	// Parse x coordinate to validate it's a number
	_, err := strconv.Atoi(xStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status: "ERROR",
			Error:  "INVALID_COORDINATES",
		})
	}

	// Format coordinate
	coord := y + xStr

	// Perform the strike
	result, err := h.GameManager.Strike(id, name, target, coord)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status: "ERROR",
			Error:  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, models.StrikeResponse{
		Status: "OK",
		Result: result,
	})
}

// @Summary Set player board
// @Description Sets a player's board configuration
// @Tags players
// @Accept json
// @Produce json
// @Param id path string true "Game ID"
// @Param name path string true "Player name"
// @Param board body models.Board true "Board configuration"
// @Success 200 {object} models.Board
// @Failure 400 {object} models.ErrorResponse
// @Router /games/{id}/players/{name}/board [post]
func (h *Handler) SetBoard(c echo.Context) error {
	// Get game ID and player name from path
	id := c.Param("id")
	name := c.Param("name")

	// Parse request body
	board := new(models.Board)
	if err := c.Bind(board); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status: "ERROR",
			Error:  "INVALID_REQUEST_BODY",
		})
	}

	// Set the board
	updatedBoard, err := h.GameManager.SetBoard(id, name, board)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status: "ERROR",
			Error:  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, updatedBoard)
}

// @Summary Get player board
// @Description Gets a player's board configuration
// @Tags players
// @Accept json
// @Produce json
// @Param id path string true "Game ID"
// @Param name path string true "Player name"
// @Success 200 {object} models.Board
// @Failure 404 {object} models.ErrorResponse
// @Router /games/{id}/players/{name}/board [get]
func (h *Handler) GetBoard(c echo.Context) error {
	// Get game ID and player name from path
	id := c.Param("id")
	name := c.Param("name")

	// Get the board
	board, err := h.GameManager.GetPlayerBoard(id, name)
	if err != nil {
		return c.JSON(http.StatusNotFound, models.ErrorResponse{
			Status: "ERROR",
			Error:  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, board)
}

// @Summary Get opponent's blind board
// @Description Gets an opponent's blind board (only hits and misses)
// @Tags players
// @Accept json
// @Produce json
// @Param id path string true "Game ID"
// @Param name path string true "Opponent name"
// @Success 200 {object} models.Board
// @Failure 404 {object} models.ErrorResponse
// @Router /games/{id}/opponent/{name}/board [get]
func (h *Handler) GetOpponentBlindBoard(c echo.Context) error {
	// Get game ID and opponent name from path
	id := c.Param("id")
	name := c.Param("name")

	// Get the blind board
	board, err := h.GameManager.GetOpponentBlindBoard(id, name)
	if err != nil {
		return c.JSON(http.StatusNotFound, models.ErrorResponse{
			Status: "ERROR",
			Error:  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, board)
}

// @Summary Get player board map
// @Description Gets an ASCII representation of a player's board
// @Tags players
// @Accept json
// @Produce text/plain
// @Param id path string true "Game ID"
// @Param name path string true "Player name"
// @Success 200 {string} string
// @Failure 404 {object} models.ErrorResponse
// @Router /games/{id}/players/{name}/board/map [get]
func (h *Handler) GetBoardMap(c echo.Context) error {
	// Get game ID and player name from path
	id := c.Param("id")
	name := c.Param("name")

	// Get the board map
	boardMap, err := h.GameManager.GetBoardMap(id, name, false)
	if err != nil {
		return c.JSON(http.StatusNotFound, models.ErrorResponse{
			Status: "ERROR",
			Error:  err.Error(),
		})
	}

	return c.String(http.StatusOK, boardMap)
}

// @Summary Get opponent board map
// @Description Gets an ASCII representation of an opponent's blind board
// @Tags players
// @Accept json
// @Produce text/plain
// @Param id path string true "Game ID"
// @Param name path string true "Opponent name"
// @Success 200 {string} string
// @Failure 404 {object} models.ErrorResponse
// @Router /games/{id}/opponent/{name}/board/map [get]
func (h *Handler) GetOpponentBoardMap(c echo.Context) error {
	// Get game ID and opponent name from path
	id := c.Param("id")
	name := c.Param("name")

	// Get the board map
	boardMap, err := h.GameManager.GetBoardMap(id, name, true)
	if err != nil {
		return c.JSON(http.StatusNotFound, models.ErrorResponse{
			Status: "ERROR",
			Error:  err.Error(),
		})
	}

	return c.String(http.StatusOK, boardMap)
}
