package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/xorduna/energywar/pkg/game"
	"github.com/xorduna/energywar/pkg/handlers"

	_ "github.com/xorduna/energywar/docs" // Import generated swagger docs
)

// useEmbedded determines whether to use embedded files or serve directly from disk
// Set to true for production deployment
const useEmbedded = true

//go:embed frontend
var frontendFS embed.FS

// @title Energy War Game API
// @version 1.0
// @description API for the Energy War Game
// @BasePath /api
func main() {
	// Create a new Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	//e.Use(middleware.CORS())

	// Create game manager
	gameManager := game.NewGameManager()

	// Create handler
	handler := handlers.NewHandler(gameManager)

	// API Group
	api := e.Group("/api")

	// Game routes
	api.POST("/games", handler.CreateGame)
	api.GET("/games/:id", handler.GetGame)

	// Player routes
	api.POST("/games/:id/players/:name/ready", handler.SetPlayerReady)
	api.POST("/games/:id/players/:name/strike", handler.Strike)
	api.POST("/games/:id/players/:name/board", handler.SetBoard)
	api.GET("/games/:id/players/:name/board", handler.GetBoard)
	api.GET("/games/:id/players/:name/board/map", handler.GetBoardMap)

	// Opponent routes
	api.GET("/games/:id/opponent/:name/board", handler.GetOpponentBlindBoard)
	api.GET("/games/:id/opponent/:name/board/map", handler.GetOpponentBoardMap)

	// Swagger documentation
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Serve frontend files
	if useEmbedded {
		// Use embedded files (for production)
		frontendSubFS, err := fs.Sub(frontendFS, "frontend")
		if err != nil {
			log.Fatal(err)
		}
		e.GET("/*", echo.WrapHandler(http.FileServer(http.FS(frontendSubFS))))
	} else {
		// Serve directly from disk (for development)
		// Check if the frontend directory exists
		if _, err := os.Stat("cmd/server/frontend"); err == nil {
			e.Static("/", "cmd/server/frontend")
		} else {
			log.Fatal("Frontend directory not found. Make sure you're running from the project root.")
		}
	}

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
