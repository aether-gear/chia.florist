package domain

type FindProductParams struct {
	Page  int
	Limit int
	ID    *string
	Name  *string
}

type ProductRepository interface {
	FindProducts(params FindProductParams) ([]Product, int, error)
}
