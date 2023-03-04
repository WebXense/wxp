package tf

import (
	"log"
	"reflect"

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
	if model == nil {
		return
	}
	if reflect.TypeOf(model).Kind() == reflect.Ptr {
		model = reflect.ValueOf(model).Elem().Interface()
	}
	if reflect.TypeOf(model).Kind() != reflect.Struct {
		log.Println("[WARNING] tf: model must be a struct")
		return
	}
	if reflect.TypeOf(model).Name() == "" {
		log.Println("[WARNING] tf: model must has a name")
		return
	}
	c.models[reflect.TypeOf(model).Name()] = model
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
	if reflect.TypeOf(model).Kind() == reflect.Ptr {
		model = reflect.ValueOf(model).Elem().Interface()
	}

	nameOfModel := reflect.TypeOf(model).Name()
	if c.generated[nameOfModel] {
		return ""
	}

	outPutStr := ""
	modelStr := "export interface " + nameOfModel + " {\n"
	numOfField := reflect.TypeOf(model).NumField()
	for i := 0; i < numOfField; i++ {
		field := reflect.TypeOf(model).Field(i)
		fieldName := field.Name
		var fieldType string
		if _, ok := c.typeMap[field.Type.Name()]; ok {
			fieldType = c.typeMap[field.Type.Name()]
		} else if field.Type.Kind() == reflect.Slice {
			if field.Type.Elem().Kind() == reflect.Struct {
				outPutStr += c.convertToInterface(reflect.New(field.Type.Elem()).Interface())
				fieldType = field.Type.Elem().Name() + "[]"
			} else {
				fieldType = c.goTypeToTsType(field.Type.Elem().Name()) + "[]"
			}
		} else if field.Type.Kind() == reflect.Struct {
			fieldType = field.Type.Name()
			outPutStr += c.convertToInterface(reflect.New(field.Type).Interface())
		} else {
			fieldType = c.goTypeToTsType(field.Type.Name())
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
