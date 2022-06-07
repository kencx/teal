package teal

type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Books       []Book `json:"books,omitempty"`
	DateAdded   string `json:"-"`
	DateUpdated string `json:"-"`
}

type CategoryService interface {
	GetCategories() ([]*Category, error)
	GetCategory(id int) (*Category, error)
	CreateCategory(b *Category) error
	UpdateCategory(id int) error
	DeleteCategory(id int) error
}
