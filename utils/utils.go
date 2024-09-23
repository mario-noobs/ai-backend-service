package utils

import (
	"fmt"
	"reflect"
)

func StructToMap(data interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	v := reflect.ValueOf(data)

	// Check if it's a struct
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("provided data is not a struct")
	}

	// Iterate through the fields of the struct
	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		value := v.Field(i)

		// Store the field name and its value as interface{}
		result[field.Name] = value.Interface()
	}

	return result, nil
}
