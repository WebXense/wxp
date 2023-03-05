package tf

import (
	"log"
	"reflect"

	"github.com/WebXense/wxp/reflect_util"
	"github.com/iancoleman/strcase"
)

func New() *converter {
	return &converter{
		models:    make(map[string]interface{}),
		typeMap:   make(map[string]string),
		generated: make(map[string]bool),
	}
}

type converter struct {
	models    map[string]interface{}
	typeMap   map[string]string
	generated map[string]bool
}

func (c *converter) Add(model interface{}) {
	if model == nil ||
		!reflect_util.IsStruct(model) ||
		reflect_util.ModelName(model) == "" {
		return
	}
	c.models[reflect.TypeOf(model).Name()] = reflect_util.EnsureNotPointer(model)
}

func (c *converter) SetupTypeMap(typeMap map[string]string) {
	c.typeMap = typeMap
}

func (c *converter) ToString() string {
	modelStr := ""

	for _, model := range c.models {
		modelStr += c.convertToInterface(model)
	}

	return modelStr
}

func (c *converter) convertToInterface(model any) string {
	model = reflect_util.EnsureNotPointer(model)
	nameOfModel := reflect_util.ModelName(model)
	if c.generated[nameOfModel] {
		return ""
	}

	outPutStr := ""
	modelStr := "export interface " + nameOfModel + " {\n"

	for i := 0; i < reflect_util.NumOfField(model); i++ {
		fieldName := c.getFieldName(model, i)
		if fieldName == "" {
			continue
		}
		var fieldType string
		if _, ok := c.typeMap[reflect_util.ModelName(model)]; ok {
			fieldType = c.typeMap[reflect_util.ModelName(model)]
		} else if reflect_util.IsSlice(model) {
			if reflect_util.SliceReflectKind(reflect_util.EnsureNotPointer(model)) == reflect.Struct {
				outPutStr += c.convertToInterface(reflect_util.ValueByIndex(model, i))
				fieldType = reflect_util.SliceTypeNameByIndex(model, i) + "[]"
			} else {
				fieldType = c.goTypeToTsType(reflect_util.TypeNameByIndex(model, i)) + "[]"
			}
		} else if reflect_util.IsStruct(model) {
			field := reflect.TypeOf(model).Field(i)
			outPutStr += c.convertToInterface(reflect.New(field.Type).Interface())
		} else {
			fieldType = c.goTypeToTsType(reflect_util.TypeNameByIndex(model, i))
		}
		modelStr += "    " + strcase.ToLowerCamel(fieldName) + ": " + fieldType + ";\n"
	}

	modelStr += "}\n\n"
	outPutStr += modelStr
	c.generated[nameOfModel] = true
	return outPutStr
}

func (c *converter) goTypeToTsType(goType string) string {
	switch goType {
	case "string":
		return "string"
	case "int":
		return "number"
	case "int64":
		return "number"
	case "uint":
		return "number"
	case "uint64":
		return "number"
	case "float64":
		return "number"
	case "float32":
		return "number"
	case "bool":
		return "boolean"
	default:
		log.Fatalln("[ERROR] tf: unknown type:", goType)
	}

	return ""
}

func (c *converter) getFieldName(model interface{}, index int) string {
	for _, tag := range []string{"json", "form", "uri"} {
		if reflect_util.HasTagByIndex(model, index, tag) {
			return reflect_util.TagValueByIndex(model, index, tag)
		}
	}
	return ""
}
