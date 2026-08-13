package main

import (
	"log"

	"github.com/TNJKL/bookmark-user-service/internal/infrastructure"
)

// @title       User API for Bookmark Management Backend
// @version     4.0.0
// @description API Swagger for bookmark-user-service.
// @BasePath    /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	//Init api
	a := infrastructure.CreateAPI()

	//Start api app
	if err := a.Start(); err != nil {
		log.Fatalf("Failed to start API server: %v", err)
	}
}
