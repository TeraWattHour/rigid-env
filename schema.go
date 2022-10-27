package env

import (
	"fmt"
	"reflect"
	"strconv"
)

func Load(target interface{}, files ...string) error {
	vars, err := loadFiles(files...)
	if err != nil {
		return err
	}

	val := reflect.ValueOf(target).Elem()

	for i := 0; i < val.NumField(); i++ {
		name := val.Type().Field(i).Name
		f := val.Field(i)
		ft := f.Type()
		k := ft.Kind()
		environmentValue := vars[name]

		isPointer := k == reflect.Ptr
		if environmentValue == "" && !isPointer {
			return fmt.Errorf("%v: missing environment value", name)
		}

		if isPointer {
			k = ft.Elem().Kind()
		}

		if err := setTargetValue(f, environmentValue, k, isPointer); err != nil {
			return fmt.Errorf("%v: %v", name, err)
		}

	}

	return nil
}

func setTargetValue(valueField reflect.Value, value string, kind reflect.Kind, isPointer bool) error {
	var v reflect.Value

	// casting to integer
	if kind == reflect.Int {
		integer, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("provided environment variable is not of type int")
		}
		v = reflect.ValueOf(&integer)
	} else if kind == reflect.String {
		v = reflect.ValueOf(&value)
	} else {
		return fmt.Errorf("casting to kind %v is unsupported", kind)
	}

	if !isPointer {
		valueField.Set(v.Elem())
	} else {
		valueField.Set(v)
	}

	return nil
}
