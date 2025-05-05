package product

type ProductService struct {
    Repo *ProductRepository
}

func NewProductService(repo *ProductRepository) *ProductService {
    return &ProductService{Repo: repo}
}

func (s *ProductService) CreateProduct(prod *Product) error {
    return s.Repo.CreateProduct(prod)
}
