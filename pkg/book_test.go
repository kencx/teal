package pkg

import "testing"

func TestCheckVal(t *testing.T) {
	b := &Book{
		Title:  "Foo Bar",
		Author: "John Doe",
		ISBN:   "12367",
	}

	if err := b.Validate(); err != nil {
		t.Error(err)
	}
}
