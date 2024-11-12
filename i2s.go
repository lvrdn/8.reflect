package main

import (
	"fmt"
	"reflect"
)

func i2s(data interface{}, out interface{}) error {

	dataValue := reflect.ValueOf(data)

	outKind := reflect.ValueOf(out).Kind()
	if outKind != reflect.Interface && outKind != reflect.Pointer {
		return fmt.Errorf("error: expected interface{} or pointer, got %v", outKind.String())
	}
	outValue := reflect.ValueOf(out).Elem()

	switch outValue.Type().Kind() {

	case reflect.Slice:
		dataSlice, ok := data.([]interface{})
		if !ok {
			return fmt.Errorf("error: expected []interface{}, got %v", dataValue.Type().String())
		}

		typeOfSlice := outValue.Type()

		newSlice := reflect.New(typeOfSlice).Elem()

		for i := 0; i < len(dataSlice); i++ {
			typeOfStructFieldInSlice := outValue.Type().Elem()

			newPtrToStruct := reflect.New(typeOfStructFieldInSlice).Interface()
			err := i2s(dataSlice[i], newPtrToStruct)
			if err != nil {
				return err
			}
			newSlice = reflect.Append(newSlice, reflect.ValueOf(newPtrToStruct).Elem())
		}
		outValue.Set(newSlice)

	case reflect.Struct:
		dataToPutInRecFunc, ok := data.(map[string]interface{})
		if !ok {
			return fmt.Errorf("error: expected struct, got %v", dataValue.Type().String())
		}

		mapKeys := dataValue.MapKeys()
		for i := 0; i < len(mapKeys); i++ {
			var typeOfStructField reflect.Type
			var newPtrToField interface{}
			var numField int
			for j := 0; j < outValue.NumField(); j++ {
				if outValue.Type().Field(j).Name == mapKeys[i].String() {
					numField = j
					typeOfStructField = outValue.Field(j).Type()
					newPtrToField = reflect.New(typeOfStructField).Interface()
					break
				}
			}
			err := i2s(dataToPutInRecFunc[mapKeys[i].String()], newPtrToField)
			if err != nil {
				return err
			}
			outValue.Field(numField).Set(reflect.ValueOf(newPtrToField).Elem())
		}

	case reflect.String:
		if _, ok := dataValue.Interface().(string); !ok {
			return fmt.Errorf("error: expected string, got %v", dataValue.Type().String())
		}
		outValue.SetString(dataValue.String())

	case reflect.Int:
		if _, ok := dataValue.Interface().(float64); !ok {
			return fmt.Errorf("error: expected float, got %v", dataValue.Type().String())
		}
		outValue.SetInt(int64(dataValue.Float()))

	case reflect.Bool:
		if _, ok := dataValue.Interface().(bool); !ok {
			return fmt.Errorf("error: expected bool, got %v", dataValue.Type().String())
		}
		outValue.SetBool(dataValue.Bool())
	}

	return nil
}
