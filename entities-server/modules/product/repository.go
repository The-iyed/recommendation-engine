package product

import "gorm.io/gorm"

type ProductRepository struct {
    DB *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
    return &ProductRepository{DB: db}
}

func (r *ProductRepository) CreateProduct(prod *Product) error {
    return r.DB.Create(prod).Error
}
