package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
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

	AppInvalid struct {
		Version string `validate:"len:lol"`
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
	MinimumWage struct {
		Currency int `validate:"min:500"`
	}
)

func TestValidateOk(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
		description string
	}{
		{
			in: Token{
				Header:    []byte("Auth"),
				Payload:   []byte(""),
				Signature: []byte(""),
			},
			expectedErr: nil,
			description: "No validate tag",
		},
		{
			in: User{
				ID:     "123456789012345678901234567890123456",
				Name:   "Dan",
				Age:    30,
				Email:  "dan@example.com",
				Role:   "admin",
				Phones: []string{"89012345678", "89023456789"},
				meta:   nil,
			},
			expectedErr: nil,
			description: "Mixed validators",
		},
		{
			in: &App{
				Version: "0.0.1",
			},
			expectedErr: nil,
			description: "struct as pointer",
		},
		{
			in: Response{
				Code: 500,
				Body: "Internal Server Error",
			},
			expectedErr: nil,
			description: "in validator",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d %s", i, tt.description), func(t *testing.T) {
			tt := tt
			t.Parallel()
			err := Validate(tt.in)
			require.Nil(t, err)
		})
	}
}

func TestValidateError(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
		description string
	}{
		{
			in: App{
				Version: "0.0.12",
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Version",
					Err:   ErrValidateLen,
				},
			},
			description: "validation len error",
		},
		{
			in: User{
				ID:     "12",
				Name:   "Dan",
				Age:    55,
				Email:  "dan@example.com.ru",
				Role:   "staff",
				Phones: []string{"12", "13"},
				meta:   nil,
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "ID",
					Err:   ErrValidateLen,
				},
				ValidationError{
					Field: "Age",
					Err:   ErrValidateMax,
				},
				ValidationError{
					Field: "Email",
					Err:   ErrValidateRegexp,
				},
				ValidationError{
					Field: "Role",
					Err:   ErrValidateIn,
				},
				ValidationError{
					Field: "Phones",
					Err:   ErrValidateLen,
				},
				ValidationError{
					Field: "Phones",
					Err:   ErrValidateLen,
				},
			},
			description: "mixed validation errors",
		},
		{
			in: Response{
				Code: 401,
				Body: "Unauthorized",
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Code",
					Err:   ErrValidateIn,
				},
			},
			description: "validation in error",
		},
		{
			in: MinimumWage{
				Currency: 100,
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Currency",
					Err:   ErrValidateMin,
				},
			},
			description: "validate min error",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d %s", i, tt.description), func(t *testing.T) {
			tt := tt
			t.Parallel()
			err := Validate(tt.in)
			require.EqualError(t, err, tt.expectedErr.Error())
		})
	}
}

func TestNonValidateError(t *testing.T) {
	var role UserRole = "user"
	type (
		Person struct {
			Name string `validate:"len:"`
		}
		InvalidtagNameInt struct {
			ID int `validate:"uniq:true"`
		}
		InvalidTagNameStr struct {
			Name string `validate:"length:3"`
		}
		InvalidType struct {
			IsExists bool `validate:"istrue:yes"`
		}
		InvalidTag struct {
			Name string `validate:""`
		}
		InvalidRegexpTag struct {
			//nolint:lll
			Emails []string `validate:"regexp:/(?!.*?\\.\\.)(?!(\\.|\\.\\.).*@.*\\.(\\.\\.|\\.)\\.*)(^[a-zA-Z][\\w.-]*[a-zA-Z]+@[a-zA-Z][\\w.-]*[a-zA-Z]+\\.[a-zA-Z][a-zA-Z.]*[a-zA-Z]$)/"`
		}
		InvalidMinTag struct {
			Currency int `validate:"min:১০০"`
		}
		InvalidMaxTag struct {
			Currency int `validate:"max:১০০"`
		}
	)

	tests := []struct {
		in          interface{}
		expectedErr error
		description string
	}{
		{
			in:          role,
			expectedErr: ErrUnsupportedType,
			description: "not a struct type",
		},
		{
			in:          nil,
			expectedErr: ErrUnsupportedType,
			description: "nil type",
		},

		{
			in: AppInvalid{
				Version: "0.0.1",
			},
			expectedErr: ErrInvalidRuleValue,
			description: "invalid rule value len",
		},

		{
			in: Person{
				Name: "Dan",
			},
			expectedErr: ErrInvalidRuleValue,
			description: "has tag no rule",
		},
		{
			in: InvalidtagNameInt{
				ID: 123123,
			},
			expectedErr: ErrInvalidRule,
			description: "unsupported validation tag for int",
		},
		{
			in: InvalidTagNameStr{
				Name: "Danil",
			},
			expectedErr: ErrInvalidRule,
			description: "unsuppordet validation tag for string",
		},
		{
			in: InvalidType{
				IsExists: true,
			},
			expectedErr: ErrUnsupportedField,
			description: "unsupported type of field",
		},
		{
			in: InvalidTag{
				Name: "Dan",
			},
			expectedErr: ErrInvalidRule,
			description: "invalid rule",
		},
		{
			in: InvalidRegexpTag{
				Emails: []string{"dan@example.com", "dan@supermail.ru"},
			},
			expectedErr: ErrInvalidRuleValue,
			description: "invalid regexp rule",
		},
		{
			in: InvalidMinTag{
				Currency: 150,
			},
			expectedErr: ErrInvalidRuleValue,
			description: "invalid min tag",
		},
		{
			in: InvalidMaxTag{
				Currency: 50,
			},
			expectedErr: ErrInvalidRuleValue,
			description: "invalid max tag",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d %s", i, tt.description), func(t *testing.T) {
			tt := tt
			t.Parallel()
			err := Validate(tt.in)
			require.ErrorIs(t, err, tt.expectedErr)
		})
	}
}
