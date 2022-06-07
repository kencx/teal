package teal

type Series struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Books       []Book `json:"books,omitempty"`
	DateAdded   string `json:"-"`
	DateUpdated string `json:"-"`
}

type SeriesService interface {
	GetAllSeries() ([]*Series, error)
	GetSeries(id int) (*Series, error)
	CreateSeries(b *Series) error
	UpdateSeries(id int) error
	DeleteSeries(id int) error
}
