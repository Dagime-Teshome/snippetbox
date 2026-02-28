package forms

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"unicode/utf8"
)

type Form struct {
	url.Values
	Errors errors
}

var EmailRX = regexp.MustCompile(
	`^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9-]+(\.[a-zA-Z0-9-]+)*$`,
)

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
			f.Errors.Add(field, fmt.Sprintf("%s This field cannot be blank", field))
		}
	}
}

func (f *Form) MaxLength(field string, limit int) {

	value := f.Get(field)
	if utf8.RuneCountInString(value) > limit {
		f.Errors.Add(field, fmt.Sprintf("%s cant be more than %d", field, limit))
	}
}
func (f *Form) MinLength(field string, limit int) {
	value := f.Get(field)
	if utf8.RuneCountInString(value) < limit {
		f.Errors.Add(field, fmt.Sprintf("%s Password is too short :password must have a min length of %v", field, limit))
	}
}

func (f *Form) MatchesPattern(field string, pattern *regexp.Regexp) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if !pattern.MatchString(value) {
		f.Errors.Add(field, "Enter a valid email address")
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
