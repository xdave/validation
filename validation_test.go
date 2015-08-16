package validation_test

import (
	"encoding/json"
	"fmt"
)

import (
	"github.com/xdave/validation"
)

type Employer struct {
	Name string `json:"name" valid:"Required"`
}

type Person struct {
	Name     string `json:"name" valid:"Required"`
	Age      int    `json:"age" valid:"Required;Min(18)"`
	Employer Employer
}

func ExampleValidationSuccess() {
	obj := Person{}
	input := `{
                "name": "John",
                "age": 35,
                "employer": {
                  "name": "Widgets, Inc."
                }
              }`
	if err := json.Unmarshal([]byte(input), &obj); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(obj)
	// Output: {John 35 {Widgets, Inc.}}
}

func ExampleValidationFailure() {
	obj := Person{}
	input := `{ "age": 17 }`
	if err := json.Unmarshal([]byte(input), &obj); err != nil {
		fmt.Println(err)
		return
	}
	if errs := validation.Validate(obj); len(errs) > 0 {
		for _, err := range errs {
			fmt.Println(err)
		}
		return
	}
	fmt.Println(obj)
	// Output:
	// validation_test.Person validation failed: `Name` Can not be empty (actual value: "")
	// validation_test.Person validation failed: `Age` Minimum is 18 (actual value: 17)
	// validation_test.Person.Employer validation failed: `Name` Can not be empty (actual value: "")
}

func ExampleValidationFailure2() {
	obj := Person{}
	input := `{ "name": "Sam", "age": 18 }`
	if err := json.Unmarshal([]byte(input), &obj); err != nil {
		fmt.Println(err)
		return
	}
	if errs := validation.Validate(obj); len(errs) > 0 {
		for _, err := range errs {
			fmt.Println(err)
		}
		return
	}
	fmt.Println(obj)
	// Output:
	// validation_test.Person.Employer validation failed: `Name` Can not be empty (actual value: "")
}
