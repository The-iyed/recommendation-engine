package router

import (
	"entities-server/modules/cloudinary"
	"entities-server/modules/product"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()
	cld, err := cloudinary.New()
	if err != nil {
		log.Fatalf("Failed to initialize Cloudinary: %v", err)
	}
	productRepo := product.NewProductRepository(db)
	productService := product.NewProductService(productRepo)
	productController := product.NewProductController(productService, cld)

	router.POST("/products", productController.CreateProduct)

	return router
}
