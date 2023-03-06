package wxp

import (
	"reflect"
)

type BaseMapper[T any] struct{}

func (m *BaseMapper[T]) Map2DTO(from interface{}) *T {
	return MapObject[T](from)
}

func (m *BaseMapper[T]) Map2DTOs(fromArray interface{}) []T {
	fromVal := reflect.ValueOf(fromArray)
	if fromVal.Kind() != reflect.Slice {
		panic("from must be a slice")
	}

	var from = make([]interface{}, fromVal.Len())
	for i := 0; i < fromVal.Len(); i++ {
		from[i] = fromVal.Index(i).Interface()
	}

	var to []T
	for _, f := range from {
		t := m.Map2DTO(f)
		to = append(to, *t)
	}
	return to
}

func MapObject[T any](from interface{}) *T {
	to := reflect.ValueOf(new(T)).Elem()

	if from == nil {
		return nil
	}
	if reflect.TypeOf(from).Kind() == reflect.Ptr {
		from = reflect.ValueOf(from).Elem().Interface()
	}

	val := reflect.ValueOf(from)
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if !field.IsZero() {
			switch field.Kind() {
			case reflect.String:
				to.Field(i).SetString(field.String())
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				to.Field(i).SetInt(field.Int())
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				to.Field(i).SetUint(field.Uint())
			case reflect.Float32, reflect.Float64:
				to.Field(i).SetFloat(field.Float())
			case reflect.Bool:
				to.Field(i).SetBool(field.Bool())
			case reflect.Slice, reflect.Array, reflect.Map, reflect.Struct, reflect.Ptr:
				to.Field(i).Set(field)
			}
		}
	}

	output := to.Interface().(T)
	return &output
}
