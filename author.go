package teal

import (
	"fmt"

	"github.com/kencx/teal/validator"
)

type Author struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (a Author) String() string {
	return fmt.Sprintf(`Author %s`, a.Name)
}

func (a *Author) Validate(v *validator.Validator) {
	v.Check(a.Name != "", "name", "value is missing")
}
