package router

import (
	"r-builder/internal/handlers"
	"r-builder/modules/graph"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	db  , _ := graph.NewNeo4jDB()

	router.GET("/recommend-product", func(c *gin.Context) {
		handlers.GetRecommendedProducts(c, db)
	})

	return router
}
