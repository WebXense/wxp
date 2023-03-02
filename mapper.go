package wxp

import (
	"fmt"
	"reflect"
)

type BaseMapper[T any] struct{}

func (m *BaseMapper[T]) Map2DTO(from interface{}) (*T, error) {
	return m.map2(from)
}

func (m *BaseMapper[T]) Map2DTOs(from []interface{}) ([]T, error) {
	var to []T
	for _, f := range from {
		t, err := m.Map2DTO(f)
		if err != nil {
			return nil, err
		}
		to = append(to, *t)
	}
	return to, nil
}

func (m *BaseMapper[T]) map2(from interface{}) (*T, error) {
	val := reflect.Indirect(reflect.ValueOf(from))
	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("from must be a struct")
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
	return to, nil
}
