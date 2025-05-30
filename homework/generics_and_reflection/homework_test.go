package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"strings"
	"testing"
)

// go test -v homework_test.go

type Person struct {
	Name    string `properties:"name"`
	Address string `properties:"address,omitempty"`
	Age     int    `properties:"age"`
	Married bool   `properties:"married"`
}

func Serialize(person Person) string {
	typ := reflect.TypeOf(person)
	val := reflect.ValueOf(person)

	resultString := strings.Builder{}
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldName, isOmitempty := getPropertyTag(field)

		value := val.Field(i).Interface()
		valToString := fmt.Sprint(value)
		if valToString == "" && isOmitempty {
			continue
		}

		resultString.WriteString(fieldName)
		resultString.WriteString("=")
		resultString.WriteString(fmt.Sprint(value))

		if i == typ.NumField()-1 {
			break
		}
		resultString.WriteString("\n")
	}

	return resultString.String()
}

func getPropertyTag(field reflect.StructField) (fieldName string, isOmitempty bool) {
	tag := field.Tag.Get("properties")
	tagFields := strings.Split(tag, ",")
	if len(tagFields) == 0 {
		return
	}
	fieldName = tagFields[0]
	for _, v := range tagFields {
		if v == "omitempty" {
			isOmitempty = true
			break
		}
	}
	return
}

func TestSerialization(t *testing.T) {
	tests := map[string]struct {
		person Person
		result string
	}{
		"test case with empty fields": {
			result: "name=\nage=0\nmarried=false",
		},
		"test case with fields": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
			},
			result: "name=John Doe\nage=30\nmarried=true",
		},
		"test case with omitempty field": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
				Address: "Paris",
			},
			result: "name=John Doe\naddress=Paris\nage=30\nmarried=true",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := Serialize(test.person)
			assert.Equal(t, test.result, result)
		})
	}
}
