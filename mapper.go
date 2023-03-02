package wxp

import (
	"reflect"
)

type BaseMapper[T any] struct{}

func (m *BaseMapper[T]) Map2DTO(from interface{}) *T {
	return MapObject(from, new(T))
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

func MapObject[T any](from interface{}, to *T) *T {
	toVal := reflect.Indirect(reflect.ValueOf(to))
	if toVal.Kind() != reflect.Struct {
		panic("to must be a struct")
	}

	fromVal := reflect.Indirect(reflect.ValueOf(from))
	if fromVal.Kind() != reflect.Struct {
		panic("from must be a struct")
	}

	for i := 0; i < toVal.NumField(); i++ {
		field := toVal.Type().Field(i)
		if val, ok := fromVal.Type().FieldByName(field.Name); ok {
			toVal.Field(i).Set(fromVal.Field(val.Index[0]))
		}
	}

	return to
}
