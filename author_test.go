package teal

import (
	"strings"
	"testing"
)

func TestValidateAuthor(t *testing.T) {
	tests := []struct {
		name   string
		author *Author
		err    string
	}{{
		name: "success",
		author: &Author{
			Name: "John Doe",
		},
		err: "",
	}, {
		name: "no name",
		author: &Author{
			Name: "",
		},
		err: "name",
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			verr := tt.author.Validate()

			if verr != nil && tt.err == "" {
				t.Fatalf("expected no err, got %v", verr)
			}

			if verr == nil && tt.err != "" {
				t.Fatalf("expected err with %q, got nil", tt.err)
			}

			if verr != nil && tt.err != "" {
				if !strings.Contains(verr.Message, tt.err) {
					t.Fatalf("got %v, want %v error", verr.Message, tt.err)
				}
			}
		})
	}

}
