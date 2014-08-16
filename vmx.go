package vmx

import (
	"fmt"
	"reflect"
)

// Returns VMX data generated from the Go value v
func Marshal(v interface{}) ([]byte, error) {
	val := reflect.ValueOf(v).Elem()

	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		tag := typeField.Tag

		fmt.Printf("Field Name: %s,\t Field Value: %v,\t Tag Value: %s\n", typeField.Name, valueField.Interface(), tag.Get("tag_name"))
	}
	return nil, nil
}

// Takes VMX data and binds it to the Go value pointed by v
func Unmarshal(data []byte, v interface{}) error {
	return nil
}
