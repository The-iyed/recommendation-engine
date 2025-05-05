package product

import (
	"entities-server/modules/cloudinary"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductController struct {
	Service    *ProductService
	Cloudinary *cloudinary.CloudinaryUploader
}

func NewProductController(service *ProductService, cld *cloudinary.CloudinaryUploader) *ProductController {
	return &ProductController{
		Service:    service,
		Cloudinary: cld,
	}
}

type CreateProductResponse struct {
	Message string
	Product Product
}

func (c *ProductController) CreateProduct(ctx *gin.Context) {
	var prod Product

	prod.Name = ctx.PostForm("name")
	prod.Description = ctx.PostForm("description")
	price, err := strconv.ParseFloat(ctx.PostForm("price"), 64)
	if err != nil || price <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid price"})
		return
	}
	prod.Price = price

	file, err := ctx.FormFile("image")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Image file is required"})
		return
	}
	fileContent, err := file.Open()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open image file"})
		return
	}
	defer fileContent.Close()

	uploadResult, err := c.Cloudinary.UploadImage(fileContent)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image to Cloudinary"})
		return
	}
	fmt.Println("Image uploaded to Cloudinary:", uploadResult)

	prod.ImagePath = uploadResult

	if err := c.Service.CreateProduct(&prod); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}

	ctx.JSON(http.StatusOK, &CreateProductResponse{
		Message: "Product created with success !",
		Product: prod,
	})
}
