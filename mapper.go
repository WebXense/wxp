package wxp

import (
	"reflect"
)

type BaseMapper[T any] struct{}

func (m *BaseMapper[T]) Map2DTO(from interface{}) *T {
	return m.map2(from)
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

func (m *BaseMapper[T]) map2(from interface{}) *T {
	val := reflect.Indirect(reflect.ValueOf(from))
	if val.Kind() != reflect.Struct {
		panic("from must be a struct")
	}

	values := make(map[string]reflect.Value)
	for i := 0; i < val.NumField(); i++ {
		values[val.Type().Field(i).Name] = val.Field(i)
	}

	to := new(T)
	toVal := reflect.Indirect(reflect.ValueOf(to))
	for i := 0; i < toVal.NumField(); i++ {
		field := toVal.Type().Field(i)
		if val, ok := values[field.Name]; ok {
			toVal.Field(i).Set(val)
		}
	}
	return to
}
