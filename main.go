package main

import (

	"github.com/b055/cards/models"

	"github.com/b055/cards/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	models.ConnectDatabase()
	r := gin.Default()

	// API v1
	v1 := r.Group("/api/v1") // versioned API is pretty important
	{
		v1.GET("decks", handlers.GetAllDecks)
		v1.GET("decks/:deck_id", handlers.GetDeckById)
		v1.POST("decks", handlers.CreateDeck)
		v1.GET("decks/:deck_id/*draw", handlers.DrawCardsInDeck)
	}

	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	r.Run()
}
