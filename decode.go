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

		sourceKey := strings.TrimSpace(parts[0])
		sourceValue := strings.TrimSpace(parts[1])
		sourceValue = strings.TrimPrefix(sourceValue, `"`)
		sourceValue = strings.TrimSuffix(sourceValue, `"`)

		wg.Add(1)
		go func(val reflect.Value, key, sourceKey, sourceValue string) {
			defer wg.Done()
			err := d.decode(val, key, sourceKey, sourceValue)
			if err != nil {
				log.Printf("Error decoding: %s\n", err)
			}
		}(val, "", sourceKey, sourceValue)
	}
	wg.Wait()

	if err := d.scanner.Err(); err != nil {
		return fmt.Errorf("Scanner error: %v", err)
	}

	return nil
}

func (d *Decoder) decode(val reflect.Value, key, sourceKey, sourceValue string) error {
	var err error
	sourceKey = strings.ToLower(sourceKey)

	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		tag := string(typeField.Tag)

		// If field does not have a tag, do not bind any value.
		if tag == "" {
			continue
		}

		destKey, _, err := parseTag(tag)
		if err != nil {
			continue
		}

		if key != "" {
			destKey = key + "." + destKey
		}

		destKey = strings.ToLower(destKey)

		//fmt.Printf("\n->%s<- has prefix ->%s<-? %t ->\n", sourceKey, destKey, strings.HasPrefix(sourceKey, destKey))

		if destKey == "-" || !strings.HasPrefix(sourceKey, destKey) || !valueField.CanSet() {
			continue
		}

		switch valueField.Kind() {
		case reflect.Struct:
			err = d.decode(valueField, destKey, sourceKey, sourceValue)

		case reflect.Array, reflect.Slice:
			err = d.decodeArray(valueField, destKey, sourceKey, sourceValue)

		case reflect.String:
			valueField.SetString(sourceValue)

		case reflect.Bool:
			var boolValue bool
			boolValue, err = strconv.ParseBool(sourceValue)
			valueField.SetBool(boolValue)

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			var intValue int64
			intValue, err = strconv.ParseInt(sourceValue, 10, valueField.Type().Bits())
			valueField.SetInt(intValue)

		case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint8:
			var uintValue uint64
			uintValue, err = strconv.ParseUint(sourceValue, 10, valueField.Type().Bits())
			valueField.SetUint(uintValue)

		default:
			err = fmt.Errorf("data type unsupported: %s", valueField.Kind())
		}
	}
	return err
}

func (d *Decoder) decodeArray(valueField reflect.Value, destKey, sourceKey, sourceValue string) error {
	fmt.Printf("Dest key => %s, ", destKey)
	fmt.Printf("Source key => %s, ", sourceKey)
	fmt.Printf("Source value => %s\n", sourceValue)

	// TODO(c4milo): I need to figure out how to grow the slice using the
	// reflection API

	length := valueField.Len()
	capacity := valueField.Cap()
	if length >= capacity {
		capacity := 2 * length
		if capacity < 4 {
			capacity = 4
		}
		newSlice := reflect.MakeSlice(valueField.Type(), length, capacity)
		reflect.Copy(newSlice, valueField)
		valueField.Set(newSlice)
	}
	valueField.SetLen(length + 1)
	destKey += strconv.Itoa(length)

	err := d.decode(valueField.Index(length), destKey, sourceKey, sourceValue)
	if err != nil {
		valueField.SetLen(length)
	}
	return err
}

func searchValue(valueField reflect.Value, key string) (uint, bool) {
	return 0, false
}
