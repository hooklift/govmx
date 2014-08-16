package vmx

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type Encoder struct {
	// Buffer where the vmx file will be written to
	buffer *bytes.Buffer
	// Maximum recursion allowed
	maxRecursion uint8
	// Current recursion level
	currentRecursion uint8
	// Parent key, used for recursion. This will allow us to set the correct
	// keys for nested structures.
	parentKey string
}

func NewEncoder(buffer *bytes.Buffer) *Encoder {
	return &Encoder{
		buffer:       buffer,
		maxRecursion: 5,
	}
}

func (e *Encoder) Encode(v interface{}) error {
	val := reflect.ValueOf(v)
	return e.encode(val)
}

func (e *Encoder) encode(val reflect.Value) error {
	// Drill into interfaces and pointers.
	// This can turn into an infinite loop given a cyclic chain,
	// but it matches the Go 1 behavior.
	for val.Kind() == reflect.Interface || val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil
		}
		val = val.Elem()
	}

	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		tag := typeField.Tag

		key, omitempty, err := parseTag(string(tag))
		if err != nil {
			return err
		}

		if key == "-" || !valueField.IsValid() ||
			(omitempty && isEmptyValue(valueField)) {
			continue
		}

		switch valueField.Kind() {
		case reflect.Struct:
			err = e.encodeStruct(valueField, key)
		case reflect.Array, reflect.Slice:
			err = e.encodeArray(valueField, key)
		default:
			if !valueField.CanSet() {
				continue
			}
			if e.parentKey != "" {
				key = e.parentKey + "." + key
			}
			e.buffer.WriteString(fmt.Sprintf("%s = \"%v\"\n", key, valueField.Interface()))
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func (e *Encoder) encodeArray(valueField reflect.Value, key string) error {
	for i := 0; i < valueField.Len(); i++ {
		e.parentKey = key + strconv.Itoa(i)

		err := e.encode(valueField.Index(i))
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *Encoder) encodeStruct(valueField reflect.Value, key string) error {
	e.currentRecursion++
	if e.currentRecursion > e.maxRecursion {
		return nil
	}

	e.parentKey += key
	err := e.encode(valueField)
	if err != nil {
		return err
	}
	e.parentKey = ""
	e.currentRecursion--
	return nil
}

func parseTag(tag string) (string, bool, error) {
	omitempty := false

	// Takes out first colon found
	parts := strings.Split(tag, ":")
	if len(parts) < 2 || parts[1] == "" {
		return "", omitempty, fmt.Errorf("Invalid tag: %s", tag)
	}

	if parts[1] == `""` {
		return "", omitempty, fmt.Errorf("Tag name is missing: %s", tag)
	}

	// Takes out double quotes
	parts2 := strings.Split(parts[1], `"`)
	if len(parts2) < 2 {
		return "", omitempty, fmt.Errorf("Tag name has to be enclosed in double quotes: %s", tag)
	}

	values := strings.Split(parts2[1], ",")
	if len(values) > 1 && values[1] == "omitempty" {
		omitempty = true

	}

	return values[0], omitempty, nil
}

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}
