package pkg

type Author struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Books       []Book `json:"books,omitempty"`
	DateAdded   string `json:"-"`
	DateUpdated string `json:"-"`
}

type AuthorService interface {
	GetAuthor(id int) (*Author, error)
	GetAllAuthors() ([]*Author, error)
	CreateAuthor(b *Author) (int, error)
	UpdateAuthor(id int) error
	DeleteAuthor(id int) error
}
