package teal

import (
	"strings"
	"testing"
)

func TestValidateBook(t *testing.T) {
	tests := []struct {
		name string
		book *Book
		err  string
	}{{
		name: "success",
		book: &Book{
			Title:  "FooBar",
			ISBN:   "100",
			Author: []string{"John Doe"},
		},
		err: "",
	}, {
		name: "no title",
		book: &Book{
			ISBN:   "100",
			Author: []string{"John Doe"},
		},
		err: "title",
	}, {
		name: "nil author",
		book: &Book{
			Title:  "Foo Bar",
			ISBN:   "100",
			Author: nil,
		},
		err: "author",
	}, {
		name: "zero length author",
		book: &Book{
			Title:  "Foo Bar",
			ISBN:   "100",
			Author: []string{},
		},
		err: "author",
	}, {
		name: "empty string author",
		book: &Book{
			Title:  "Foo Bar",
			ISBN:   "100",
			Author: []string{""},
		},
		err: "author",
	}, {
		name: "no isbn",
		book: &Book{
			Title:  "Foo Bar",
			Author: []string{"John Doe"},
		},
		err: "isbn",
	}, {
		name: "isbn regex fail",
		book: &Book{
			Title:  "Foo Bar",
			ISBN:   "abc",
			Author: []string{"John Doe"},
		},
		err: "isbn",
	}, {
		name: "multiple errors",
		book: &Book{
			Title:  "Foo Bar",
			ISBN:   "abc",
			Author: []string{""},
		},
		err: "author, isbn",
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			verrs := tt.book.Validate()

			if len(verrs) != 0 && tt.err == "" {
				t.Fatalf("expected no err, got %d errors", len(verrs))
			}

			if len(verrs) == 0 && tt.err != "" {
				t.Fatalf("expected err with %q, got 0 errs", tt.err)
			}

			if len(verrs) != 0 && tt.err != "" {
				strs := strings.Split(tt.err, ", ")

				// TODO probably a better way to do this
				for _, v := range verrs {
					ok := false
					for _, s := range strs {
						if strings.Contains(v.Message, s) {
							ok = true
						}
					}
					if !ok {
						t.Fatalf("got %v, want %v err", v, tt.err)
					}
				}
			}
		})
	}
}
