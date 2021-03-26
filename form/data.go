package form

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

// PhoneRX represents phone number maching pattern
var PhoneRX = regexp.MustCompile("(^\\+[0-9]{2}|^\\+[0-9]{2}\\(0\\)|^\\(\\+[0-9]{2}\\)\\(0\\)|^00[0-9]{2}|^0)([0-9]{9}$|[0-9\\-\\s]{10}$)")

// EmailRX represents email address maching pattern
var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// Input represents form input values and validations
type Input struct {
	Values  url.Values
	VErrors ValidationErrors
	CSRF    string
	Cid     uint
}

// MinLength checks if a given minium length is satisfied
func (inVal *Input) MinLength(field string, d int) {
	value := inVal.Values.Get(field)
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) < d {
		inVal.VErrors.Add(field, fmt.Sprintf("This field is too short (minimum is %d characters)", d))
	}
}

func (inVal *Input) Date(field string) {
	value := inVal.Values.Get(field)

	if value == "" {
		return
	} else {
		p := strings.Split(value, " ")
		fmt.Println("mine")
		fmt.Println(p[0])
		fmt.Println("mine")
		fmt.Println("next")
		fmt.Println(p[1])
		fmt.Println("next")
		_, err := time.Parse("01/02/2006", p[0])
		if err != nil {
			inVal.VErrors.Add(field, fmt.Sprintf("Invalid date"))
		}
	}
}

// Required checks if list of provided form input fields have values
func (inVal *Input) ValidateRequiredFields(fields ...string) {
	for _, f := range fields {
		value := inVal.Values.Get(f)
		fmt.Println(value)
		if value == "" {
			fmt.Println("empty")
			fmt.Println(f)
			inVal.VErrors.Add(f, "This field is required field")
		}
	}
}

//checks if value is number
func (inVal *Input) ValidateFieldsInteger(fields ...string) {
	for _, f := range fields {
		value := inVal.Values.Get(f)

		_, err := strconv.Atoi(value)
		if err != nil {
			fmt.Println("alpabet")
			inVal.VErrors.Add(f, "This field must be a number")
		}
	}
}

// checks if negative
func (inVal *Input) ValidateFieldsRange(fields ...string) {
	for _, f := range fields {
		value := inVal.Values.Get(f)
		fmt.Println("not")
		val, err := strconv.Atoi(value)
		if err == nil && val < 0 {

			fmt.Println("negative")
			inVal.VErrors.Add(f, "This field must be positive number")
		}
	}
}

// func (inVal *Input) ValidateFieldFile(fields string) {
//       w, fh, er := r.FormFile(fields)

// 	value := inVal.Values.Get(fields)
// 	fmt.Println("not")
// 	val, err := strconv.Atoi(value)
// 	if err == nil && val < 0 {

// 		fmt.Println("empty")
// 		inVal.VErrors.Add(fields, "This field must be positive number")

// 	}
//}

//////discount range0 to 100
// func (inVal *Input) ValidatediscountRange(field string) {

// 	value := inVal.Values.Get(field)
// 	fmt.Println("not")
// 	val, err := strconv.Atoi(value)
// 	if err == nil && val > 100 {

// 		fmt.Println("empty")
// 		inVal.VErrors.Add(field, "This field must be less than 100")

// 	}
// }

// MatchesPattern checks if a given input form field matchs a given pattern
func (inVal *Input) MatchesPattern(field string, pattern *regexp.Regexp) {
	value := inVal.Values.Get(field)
	if value == "" {
		return
	}
	if !pattern.MatchString(value) {
		inVal.VErrors.Add(field, "The value entered is invalid")
	}
}

// PasswordMatches checks if Password and Confirm Password fields match
func (inVal *Input) PasswordMatches(password string, confPassword string) {
	pwd := inVal.Values.Get(password)
	confPwd := inVal.Values.Get(confPassword)
	fmt.Println("first password not match")
	fmt.Println(pwd)
	fmt.Println("first password not match")
	fmt.Println(confPwd)
	fmt.Println("first password not match")
	if pwd == "" || confPwd == "" {
		fmt.Println("third password not match")
		return
	}

	if pwd != confPwd {
		fmt.Println("secondpassword not match")
		inVal.VErrors.Add(password, "The Password and Confirm Password values did not match")
		inVal.VErrors.Add(confPassword, "The Password and Confirm Password values did not match")
	}
}

// Valid checks if any form input validation has failed or not
func (inVal *Input) IsValid() bool {
	return len(inVal.VErrors) == 0
}
