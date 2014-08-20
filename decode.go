package vmx

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

type Decoder struct {
	// Scanner to read file line by line
	scanner *bufio.Scanner
}

func NewDecoder(reader io.Reader) *Decoder {
	return &Decoder{
		scanner: bufio.NewScanner(reader),
	}
}

func (d *Decoder) Decode(v interface{}) error {
	val := reflect.ValueOf(v)

	if val.Kind() != reflect.Ptr {
		return errors.New("non-pointer value passed to Decode")
	}

	if val.IsNil() {
		return fmt.Errorf("nil value passed to Decode: %v", reflect.TypeOf(val))
	}

	// Gets setteable value
	val = val.Elem()

	if !val.CanAddr() {
		return errors.New("destination struct must be addressable")
	}

	// Starts scanning the text file
	var wg sync.WaitGroup
	for d.scanner.Scan() {
		line := d.scanner.Text()

		// Ignore comments and empty lines
		if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" {
			continue
		}

		parts := strings.Split(line, "=")
		if len(parts) != 2 {
			return fmt.Errorf("Invalid line: %s ", line)
		}

		targetKey := strings.TrimSpace(parts[0])
		targetValue := strings.TrimSpace(parts[1])
		targetValue = strings.TrimPrefix(targetValue, `"`)
		targetValue = strings.TrimSuffix(targetValue, `"`)

		wg.Add(1)
		go func(val reflect.Value, key, targetKey, targetValue string) {
			defer wg.Done()
			err := d.decode(val, key, targetKey, targetValue)
			if err != nil {
				log.Printf("Error decoding: %s\n", err)
			}
		}(val, "", targetKey, targetValue)
	}
	wg.Wait()

	if err := d.scanner.Err(); err != nil {
		return fmt.Errorf("Scanner error: %v", err)
	}

	return nil
}

func (d *Decoder) decode(val reflect.Value, key, targetKey, targetValue string) error {
	var err error
	targetKey = strings.ToLower(targetKey)

	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		tag := string(typeField.Tag)

		// If field does not have a tag, do not bind any value.
		if tag == "" {
			continue
		}

		valKey, _, err := parseTag(tag)
		if err != nil {
			continue
		}

		if key != "" {
			valKey = key + "." + valKey
		}

		valKey = strings.ToLower(valKey)

		//fmt.Printf("->%s<- has prefix ->%s<-? %t -> ", targetKey, valKey, strings.HasPrefix(targetKey, valKey))
		//fmt.Printf("can set value? %t\n", valueField.CanSet())
		if valKey == "-" || !strings.HasPrefix(targetKey, valKey) || !valueField.CanSet() {
			continue
		}

		switch valueField.Kind() {
		case reflect.Struct:
			err = d.decode(valueField, valKey, targetKey, targetValue)

		case reflect.Array, reflect.Slice:
			err = d.decodeArray(valueField, valKey, targetKey, targetValue)

		case reflect.String:
			valueField.SetString(targetValue)

		case reflect.Bool:
			var boolValue bool
			boolValue, err = strconv.ParseBool(targetValue)
			valueField.SetBool(boolValue)

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			var intValue int64
			intValue, err = strconv.ParseInt(targetValue, 10, valueField.Type().Bits())
			valueField.SetInt(intValue)

		case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint8:
			var uintValue uint64
			uintValue, err = strconv.ParseUint(targetValue, 10, valueField.Type().Bits())
			valueField.SetUint(uintValue)

		default:
			err = fmt.Errorf("data type unsupported: %s", valueField.Kind())
		}
	}
	return err
}

func (d *Decoder) decodeArray(valueField reflect.Value, key, targetKey, targetValue string) error {
	for i := 0; i < valueField.Len(); i++ {
		indexedKey := key + strconv.Itoa(i)
		fmt.Printf("indexedKey -> %s\n", indexedKey)
		err := d.decode(valueField.Index(i), indexedKey, targetKey, targetValue)
		if err != nil {
			return err
		}
	}
	return nil
}
