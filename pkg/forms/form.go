package form

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

var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9]")

func New(data url.Values) *Form {

	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

//Veriy not blank fields
func (f *Form) Required(fields ...string) {

	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank")
		}
	}
}

//Verify appropriate length field

func (f *Form) MaxLength(field string, d int) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) > d {
		f.Errors.Add(field, fmt.Sprintf("Field too long, maximum is %d", d))
	}
}

//Verify specific value for the field
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
	f.Errors.Add(field, "Field value is invalid!")
}

//Verify if there are errors
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

func (f *Form) MinLength(field string, d int) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) < d {
		f.Errors.Add(field, fmt.Sprintf("This field is too short (minimum is %d)", d))
	}
}

func (form *Form) MatchesPattern(field string, pattern *regexp.Regexp) {

	value := form.Get(field)
	if value == "" {
		return
	}

	if !pattern.MatchString(value) {
		form.Errors.Add(field, "This field is invalid")
	}
}
