package cloudinary

import (
	"context"
	"entities-server/config"
	"fmt"
	"io"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type CloudinaryUploader struct {
	client *cloudinary.Cloudinary
}

func New() (*CloudinaryUploader, error) {
	cfg := config.LoadConfig()
	fmt.Println("Cloudinary Cloud Name:", cfg.CLOUDINARY_CLOUD_NAME)
	fmt.Println("Cloudinary API Key:", cfg.CLOUDINARY_API_KEY)
	fmt.Println("Cloudinary API Secret:", cfg.CLOUDINARY_API_SECRET)

	cld, err := cloudinary.NewFromParams(
		cfg.CLOUDINARY_CLOUD_NAME,
		cfg.CLOUDINARY_API_KEY,
		cfg.CLOUDINARY_API_SECRET,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Cloudinary client: %v", err)
	}

	return &CloudinaryUploader{client: cld}, nil
}

func (c *CloudinaryUploader) UploadImage(file io.Reader) (string, error) {
	ctx := context.Background()

	uploadResult, err := c.client.Upload.Upload(ctx, file, uploader.UploadParams{})
	if err != nil {
		return "", fmt.Errorf("failed to upload image: %v", err)
	}
	fmt.Println(uploadResult)
	return uploadResult.SecureURL, nil
}
