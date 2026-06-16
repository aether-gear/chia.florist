package repository

type FindShopsParams struct {
	Page  int
	Limit int
	ID    *string
	Name  *string
}
