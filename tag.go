package teal

type Tag struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Books       []Book `json:"books,omitempty"`
	DateAdded   string `json:"-"`
	DateUpdated string `json:"-"`
}

type TagService interface {
	GetTags() ([]*Tag, error)
	GetTag(id int) (*Tag, error)
	CreateTag(b *Tag) error
	UpdateTag(id int) error
	DeleteTag(id int) error
}
