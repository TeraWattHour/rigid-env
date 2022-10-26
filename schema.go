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
	setValue := func(valueField reflect.Value, value string, kind reflect.Kind, isPointer bool) error {
		var v reflect.Value
		var vp reflect.Value

		// casting to integer
		if kind == reflect.Int {
			integer, err := strconv.Atoi(value)
			if err != nil {
				return fmt.Errorf("provided environment variable is not of type int")
			}
			v = reflect.ValueOf(integer)
			vp = reflect.ValueOf(&integer)
		} else if kind == reflect.String {
			v = reflect.ValueOf(value)
			vp = reflect.ValueOf(&value)
		}

		if !isPointer {
			valueField.Set(v)
		} else {
			valueField.Set(vp)
		}

		return nil
	}
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
			if err := setValue(f, environmentValue, k, true); err != nil {
				return fmt.Errorf("%v: %v", name, err)
			}
		} else {
			if err := setValue(f, environmentValue, k, false); err != nil {
				return fmt.Errorf("%v: %v", name, err)
			}
		}

	}

	return nil
}
