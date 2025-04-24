package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/xorduna/energywar/pkg/game"
	"github.com/xorduna/energywar/pkg/handlers"

	_ "github.com/xorduna/energywar/docs" // Import generated swagger docs
)

//go:embed frontend
var frontendFS embed.FS

// @title Energy War Game API
// @version 1.0
// @description API for the Energy War Game
// @host localhost:8080
// @BasePath /
func main() {
	// Create a new Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Create game manager
	gameManager := game.NewGameManager()

	// Create handler
	handler := handlers.NewHandler(gameManager)

	// Routes
	// Game routes
	e.POST("/games", handler.CreateGame)
	e.GET("/games/:id", handler.GetGame)

	// Player routes
	e.POST("/games/:id/players/:name/ready", handler.SetPlayerReady)
	e.POST("/games/:id/players/:name/strike", handler.Strike)
	e.POST("/games/:id/players/:name/board", handler.SetBoard)
	e.GET("/games/:id/players/:name/board", handler.GetBoard)
	e.GET("/games/:id/players/:name/board/map", handler.GetBoardMap)

	// Opponent routes
	e.GET("/games/:id/opponent/:name/board", handler.GetOpponentBlindBoard)
	e.GET("/games/:id/opponent/:name/board/map", handler.GetOpponentBoardMap)

	// Swagger documentation
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Serve frontend files
	frontendSubFS, err := fs.Sub(frontendFS, "frontend")
	if err != nil {
		log.Fatal(err)
	}
	e.GET("/*", echo.WrapHandler(http.FileServer(http.FS(frontendSubFS))))

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
