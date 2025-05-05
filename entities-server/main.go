package main

import (
	"entities-server/config"
	"entities-server/database"

	"entities-server/router"
	"log"
)

func main() {
	cfg := config.LoadConfig()
	db := database.Connect()

	r := router.SetupRouter(db)
	r.Static("/uploads", "./uploads")

	if err := r.Run(":" + cfg.Port); err != nil {

		log.Fatal("Server failed to start")
	}
}
