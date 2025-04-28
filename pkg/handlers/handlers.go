package handlers

import (
	"errors"
	"net/http"
	"strconv"

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

// validatePlayerToken checks if the provided token matches the player's token
func (h *Handler) validatePlayerToken(gameID, playerName, token string) error {
	// Get the game
	game, err := h.GameManager.GetGame(gameID)
	if err != nil {
		return errors.New("game not found")
	}

	// Check if the player exists
	playerInfo, exists := game.Players[playerName]
	if !exists {
		return errors.New("player not found")
	}

	// Check if the token matches
	if playerInfo.Token != token {
		return errors.New("invalid token")
	}

	return nil
}

// @Summary Create a new game
// @Description Creates a new game with the specified parameters
// @Tags games
// @Accept json
// @Produce json
// @Param size query int false "Board size (5-20)" default(10)
// @Param capacity query int false "Required capacity" default(1000)
// @Success 200 {object} models.Game
// @Failure 400 {object} models.ErrorResponse
// @Router /games [post]
func (h *Handler) CreateGame(c echo.Context) error {
	// Parse query parameters
	sizeStr := c.QueryParam("size")
	capacityStr := c.QueryParam("capacity")
	publicStr := c.QueryParam("public")

	// Default values
	size := 10
	capacity := 1000
	public := false

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

	// Parse public
	if publicStr != "" {
		var err error
		public, err = strconv.ParseBool(publicStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Status: "ERROR",
				Error:  "INVALID_PARAMETERS",
			})
		}
	}

	// Create the game
	gameObj, err := h.GameManager.CreateGame(size, capacity, public)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status: "ERROR",
			Error:  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, gameObj)
}

// JoinGameResponse represents the response for joining a game
type JoinGameResponse struct {
	Token string `json:"token"`
}

// @Summary Join a game
// @Description Allows a player to join an existing game
// @Tags games
// @Accept json
// @Produce json
// @Param id path string true "Game ID"
// @Param player query string true "Player name"
// @Success 200 {object} JoinGameResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /games/{id}/join [post]
func (h *Handler) JoinGame(c echo.Context) error {
	// Get game ID from path
	id := c.Param("id")

	// Get player name from query
	playerName := c.QueryParam("player")

	if playerName == "" {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status: "ERROR",
			Error:  "INVALID_PARAMETERS",
		})
	}

	// Join the game
	token, err := h.GameManager.JoinGame(id, playerName)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status: "ERROR",
			Error:  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, JoinGameResponse{
		Token: token,
	})
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

	// Create a copy of the game object with tokens hidden
	limitedGameObj := &models.Game{
		ID:     gameObj.ID,
		Status: gameObj.Status,
		Turn:   gameObj.Turn,
		Winner: gameObj.Winner,
		Players: func() map[string]models.PlayerInfo {
			limitedPlayers := make(map[string]models.PlayerInfo)
			for name, player := range gameObj.Players {
				limitedPlayers[name] = models.PlayerInfo{
					Ready:         player.Ready,
					TotalCapacity: player.TotalCapacity,
					Capacity:      player.Capacity,
				}
			}
			return limitedPlayers
		}(),
	}

	return c.JSON(http.StatusOK, limitedGameObj)
}

// @Summary Set player ready
// @Description Sets a player as ready to start the game
// @Tags players
// @Accept json
// @Produce json
// @Param id path string true "Game ID"
// @Param name path string true "Player name"
// @Param token query string true "Player token"
// @Success 200 {object} models.ReadyResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Router /games/{id}/players/{name}/ready [post]
func (h *Handler) SetPlayerReady(c echo.Context) error {
	// Get game ID and player name from path
	id := c.Param("id")
	name := c.Param("name")

	// Get token from query parameter
	token := c.QueryParam("token")
	if token == "" {
		return c.JSON(http.StatusForbidden, models.ErrorResponse{
			Status: "ERROR",
			Error:  "MISSING_TOKEN",
		})
	}

	// Validate the token
	if err := h.validatePlayerToken(id, name, token); err != nil {
		return c.JSON(http.StatusForbidden, models.ErrorResponse{
			Status: "ERROR",
			Error:  "INVALID_TOKEN",
		})
	}

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
// @Param token query string true "Player token"
// @Param target query string true "Target player name"
// @Param y query string true "Y coordinate (A-Z)"
// @Param x query int true "X coordinate (1-size)"
// @Success 200 {object} models.StrikeResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Router /games/{id}/players/{name}/strike [post]
func (h *Handler) Strike(c echo.Context) error {
	// Get game ID and player name from path
	id := c.Param("id")
	name := c.Param("name")

	// Get token from query parameter
	token := c.QueryParam("token")
	if token == "" {
		return c.JSON(http.StatusForbidden, models.ErrorResponse{
			Status: "ERROR",
			Error:  "MISSING_TOKEN",
		})
	}

	// Validate the token
	if err := h.validatePlayerToken(id, name, token); err != nil {
		return c.JSON(http.StatusForbidden, models.ErrorResponse{
			Status: "ERROR",
			Error:  "INVALID_TOKEN",
		})
	}

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
// @Param token query string true "Player token"
// @Param board body models.Board true "Board configuration"
// @Success 200 {object} models.Board
// @Failure 400 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Router /games/{id}/players/{name}/board [post]
func (h *Handler) SetBoard(c echo.Context) error {
	// Get game ID and player name from path
	id := c.Param("id")
	name := c.Param("name")

	// Get token from query parameter
	token := c.QueryParam("token")
	if token == "" {
		return c.JSON(http.StatusForbidden, models.ErrorResponse{
			Status: "ERROR",
			Error:  "MISSING_TOKEN",
		})
	}

	// Validate the token
	if err := h.validatePlayerToken(id, name, token); err != nil {
		return c.JSON(http.StatusForbidden, models.ErrorResponse{
			Status: "ERROR",
			Error:  "INVALID_TOKEN",
		})
	}

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
// @Param token query string true "Player token"
// @Success 200 {object} models.Board
// @Failure 404 {object} models.ErrorResponse
// @Router /games/{id}/players/{name}/board [get]
func (h *Handler) GetBoard(c echo.Context) error {
	// Get game ID and player name from path
	id := c.Param("id")
	name := c.Param("name")

	// Get token from query parameter
	token := c.QueryParam("token")
	// Validate the token
	if err := h.validatePlayerToken(id, name, token); err != nil {
		return c.JSON(http.StatusForbidden, models.ErrorResponse{
			Status: "ERROR",
			Error:  "INVALID_TOKEN",
		})
	}

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
// @Param token query string true "Player token"
// @Success 200 {string} string
// @Failure 403 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /games/{id}/players/{name}/board/map [get]
func (h *Handler) GetBoardMap(c echo.Context) error {
	// Get game ID and player name from path
	id := c.Param("id")
	name := c.Param("name")

	// Get token from query parameter
	token := c.QueryParam("token")
	if token == "" {
		return c.JSON(http.StatusForbidden, models.ErrorResponse{
			Status: "ERROR",
			Error:  "MISSING_TOKEN",
		})
	}

	// Validate the token
	if err := h.validatePlayerToken(id, name, token); err != nil {
		return c.JSON(http.StatusForbidden, models.ErrorResponse{
			Status: "ERROR",
			Error:  "INVALID_TOKEN",
		})
	}

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

// @Summary Get game status
// @Description Gets the limited status of a game
// @Tags games
// @Accept json
// @Produce json
// @Param id path string true "Game ID"
// @Success 200 {object} models.Game
// @Failure 404 {object} models.ErrorResponse
// @Router /games/{id}/status [get]
func (h *Handler) GetGameStatus(c echo.Context) error {
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

	// Create a limited view of the game
	limitedGameObj := &models.Game{
		ID:     gameObj.ID,
		Status: gameObj.Status,
		Turn:   gameObj.Turn,
		Winner: gameObj.Winner,
		Players: func() map[string]models.PlayerInfo {
			limitedPlayers := make(map[string]models.PlayerInfo)
			for name, player := range gameObj.Players {
				limitedPlayers[name] = models.PlayerInfo{
					Ready:         player.Ready,
					TotalCapacity: player.TotalCapacity,
					Capacity:      player.Capacity,
				}
			}
			return limitedPlayers
		}(),
	}

	return c.JSON(http.StatusOK, limitedGameObj)
}
