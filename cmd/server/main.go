package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/xorduna/energywar/pkg/config"
	"github.com/xorduna/energywar/pkg/game"
	"github.com/xorduna/energywar/pkg/handlers"
	"github.com/xorduna/energywar/pkg/models"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

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

// openDatabase opens a database connection based on the configuration
func openDatabase(cfg *config.Config) (*gorm.DB, error) {
	uri := strings.TrimSpace(cfg.Database.URI)

	var dialector gorm.Dialector
	switch {
	case uri == ":memory:" || uri == "file::memory:?cache=shared" || strings.HasPrefix(uri, "file:"):
		// Handle SQLite cases
		if uri == ":memory:" || uri == "file::memory:?cache=shared" {
			dialector = sqlite.Open("file::memory:?cache=shared")
		} else {
			// Ensure file paths are handled correctly
			if !strings.HasPrefix(uri, "file:") && !strings.HasSuffix(uri, ".db") {
				uri = "file:" + uri
			}
			if !strings.HasSuffix(uri, ".db") {
				uri += ".db"
			}
			dialector = sqlite.Open(uri)
		}
	case strings.HasPrefix(uri, "postgres://"):
		// Handle PostgreSQL
		dialector = postgres.Open(uri)
	default:
		// Default to SQLite file
		dialector = sqlite.Open("file:game.db")
	}

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %v", err)
	}

	fmt.Printf("Connected to database: %s\n", uri)

	return db, nil
}

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	// Create a new Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	//e.Use(middleware.CORS())

	// Open database
	db, err := openDatabase(cfg)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Auto migrate database models
	db.AutoMigrate(models.Game{})

	// Create game manager
	gameManager := game.NewGameManager(db)

	// Create handler
	handler := handlers.NewHandler(gameManager)

	// API Group
	api := e.Group("/api")

	// Game routes
	api.POST("/games", handler.CreateGame)
	api.GET("/games/:id", handler.GetGame)
	api.GET("/games/:id/status", handler.GetGameStatus)
	api.POST("/games/:id/join", handler.JoinGame)

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
	e.Logger.Fatal(e.Start(cfg.Server.Port))
}
