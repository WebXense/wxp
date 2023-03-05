package reflect_util

import "reflect"

func ReflectKind(model interface{}) reflect.Kind {
	return reflect.TypeOf(model).Kind()
}

func SliceReflectKind(model interface{}) reflect.Kind {
	return reflect.TypeOf(model).Elem().Kind()
}

func IsPointer(model interface{}) bool {
	return reflect.TypeOf(model).Kind() == reflect.Ptr
}

func IsStruct(model interface{}) bool {
	return reflect.TypeOf(model).Kind() == reflect.Struct
}

func IsSlice(model interface{}) bool {
	return reflect.TypeOf(model).Kind() == reflect.Slice
}

func EnsureNotPointer(model interface{}) interface{} {
	if IsPointer(model) {
		return reflect.ValueOf(model).Elem().Interface()
	}
	return model
}

func ModelName(model interface{}) string {
	return reflect.TypeOf(EnsureNotPointer(model)).Name()
}

func NumOfField(model interface{}) int {
	return reflect.TypeOf(EnsureNotPointer(model)).NumField()
}

func Fields(model interface{}) map[string]string {
	fields := make(map[string]string)
	for i := 0; i < NumOfField(EnsureNotPointer(model)); i++ {
		field := reflect.TypeOf(EnsureNotPointer(model)).Field(i)
		fields[field.Name] = field.Type.Name()
	}
	return fields
}

func FieldNameByIndex(model interface{}, index int) string {
	return reflect.TypeOf(EnsureNotPointer(model)).Field(index).Name
}

func ValueByIndex(model interface{}, index int) interface{} {
	return reflect.ValueOf(EnsureNotPointer(model)).Field(index).Interface()
}

func TypeNameByIndex(model interface{}, index int) string {
	return reflect.TypeOf(EnsureNotPointer(model)).Field(index).Type.Name()
}

func HasTagByIndex(model interface{}, index int, tag string) bool {
	field := reflect.TypeOf(EnsureNotPointer(model)).Field(index)
	return field.Tag.Get(tag) != ""
}

func TagValueByIndex(model interface{}, index int, tag string) string {
	field := reflect.TypeOf(EnsureNotPointer(model)).Field(index)
	return field.Tag.Get(tag)
}

func SliceTypeNameByIndex(model interface{}, index int) string {
	return reflect.TypeOf(model).Field(index).Type.Elem().Name()
}
