package teal

type Author struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (a *Author) Validate() *ValidationError {
	if a.Name == "" {
		return NewValidationError("name", "value is missing")
	}
	return nil
}
