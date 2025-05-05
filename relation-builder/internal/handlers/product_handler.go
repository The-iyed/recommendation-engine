package handlers

import (
	"context"
	"log"
	"net/http"
	"r-builder/internal/services"
	"r-builder/modules/graph"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetRecommendedProducts(c *gin.Context, db *graph.Neo4jDB) {

	baseProductIDStr := c.DefaultQuery("base_product_id", "")

	baseProductID, err := strconv.ParseFloat(baseProductIDStr, 64)
	if err != nil {
		log.Println("Invalid base_product_id:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid base_product_id"})
		return
	}

	products, err := services.RecommendProducts(context.Background(), db, baseProductID)
	if err != nil {
		log.Println("Error getting recommendations:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching recommended products"})
		return
	}

	c.JSON(http.StatusOK, products)
}
