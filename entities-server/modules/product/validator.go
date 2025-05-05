package product

import "github.com/go-playground/validator/v10"

var validate *validator.Validate

func ValidateProduct(prod *Product) error {
	validate = validator.New()
	return validate.Struct(prod)
}
