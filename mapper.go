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
			fieldName := val.Type().Field(i).Name
			_, ok := to.Type().FieldByName(fieldName)
			if ok {
				switch field.Kind() {
				case reflect.String:
					to.FieldByName(fieldName).SetString(field.String())
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					to.FieldByName(fieldName).SetInt(field.Int())
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					to.FieldByName(fieldName).SetUint(field.Uint())
				case reflect.Float32, reflect.Float64:
					to.FieldByName(fieldName).SetFloat(field.Float())
				case reflect.Bool:
					to.FieldByName(fieldName).SetBool(field.Bool())
				case reflect.Slice, reflect.Array, reflect.Struct, reflect.Map, reflect.Ptr:
					to.FieldByName(fieldName).Set(field)
				}
			}
		}
	}

	// handle gorm.Model
	_, ok := reflect.TypeOf(from).FieldByName("Model")
	if ok {
		for _, fieldName := range []string{"ID", "CreatedAt", "UpdatedAt", "DeletedAt"} {
			_, ok := reflect.ValueOf(from).FieldByName("Model").Type().FieldByName(fieldName)
			if !ok {
				continue
			}
			field := reflect.ValueOf(from).FieldByName("Model").FieldByName(fieldName)
			if !field.IsZero() {
				to.FieldByName(fieldName).Set(field)
			}
		}
	}

	output := to.Interface().(T)
	return &output
}
