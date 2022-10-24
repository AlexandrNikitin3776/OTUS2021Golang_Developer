package hw09structvalidator

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name        string
		in          interface{}
		expectedErr error
	}{
		{
			"unsupported type",
			42,
			UnsupportedInputType,
		},
		{
			"validate user",
			User{"id", "user", 25, "example@otus.ru", "user", []string{"1", "2"}, []byte("json")},
			ValidationErrors{
				{"ID", InvalidStringLen},
				{"Role", InvalidStringIn},
				{"Phones", InvalidStringLen},
			},
		},
		{
			"validate app",
			App{"0"},
			ValidationErrors{
				{"Version", InvalidStringLen},
			},
		},
		{
			"validate response",
			Response{155, "322"},
			ValidationErrors{
				{"Code", InvalidIntIn},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Parallel()
			err := Validate(tt.in)
			require.Equal(t, tt.expectedErr, err, "errors should be equal")
			_ = tt
		})
	}
}
