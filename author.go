package teal

type Author struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Books Books  `json:"books,omitempty"`
}

type Authors []Author
