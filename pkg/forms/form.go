package forms

import (
	"fmt"
	"net/url"
	"strings"
	"unicode/utf8"
)

type Form struct {
	url.Values
	Errors errors
}

func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

func (f *Form) Required(fields ...string) {

	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, fmt.Sprintf("%s field can not be empty", field))
		}
	}
}

func (f *Form) MaxLength(field string, limit int) {

	value := f.Get(field)
	if utf8.RuneCountInString(value) > limit {
		f.Errors.Add(field, fmt.Sprintf("%s cant be more than %d", field, limit))
	}
}

func (f *Form) PermittedValues(field string, opts ...string) {
	value := f.Get(field)
	if value == "" {
		return
	}
	for _, opt := range opts {
		if value == opt {
			return
		}
	}
	f.Errors.Add(field, fmt.Sprintf("%s has invalid value", field))

}

func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}
