package typescript

import (
	"reflect"
	"strings"

	"github.com/WebXense/wxp/logger"
)

func NewModelConverter() *modelConverter {
	return &modelConverter{
		models:    make(map[string]interface{}),
		typeMap:   make(map[string]string),
		generated: make(map[string]bool),
	}
}

type modelConverter struct {
	models    map[string]interface{}
	typeMap   map[string]string
	generated map[string]bool
}

func (c *modelConverter) Add(model interface{}) {
	if model == nil {
		return
	}
	if reflect.TypeOf(model).Kind() == reflect.Ptr {
		model = reflect.ValueOf(model).Elem().Interface()
	}
	if reflect.TypeOf(model).Kind() != reflect.Struct {
		logger.Warn("model must be a struct", model)
		return
	}
	if reflect.TypeOf(model).Name() == "" {
		logger.Warn("model must has a name", model)
		return
	}
	c.models[reflect.TypeOf(model).Name()] = model
}

func (c *modelConverter) SetupTypeMap(typeMap map[string]string) {
	c.typeMap = typeMap
}

func (c *modelConverter) ToString() string {
	modelStr := ""

	for _, model := range c.models {
		modelStr += c.convertToInterface(model)
	}

	return modelStr
}

func (c *modelConverter) convertToInterface(model any) string {
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
		fieldName := c.getFieldName(model, i)
		if fieldName == "" {
			continue
		}
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
		modelStr += "    " + fieldName + ": " + fieldType + ";\n"
	}

	modelStr += "}\n\n"
	outPutStr += modelStr
	c.generated[nameOfModel] = true
	return outPutStr
}

func (c *modelConverter) goTypeToTsType(goType string) string {
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
		logger.Err("unknown type", goType)
	}

	return ""
}

func (c *modelConverter) getFieldName(model interface{}, index int) string {
	for _, tag := range []string{"json", "form", "uri"} {
		t := reflect.TypeOf(model).Field(index).Tag.Get(tag)
		if t != "" {
			var canOmit bool
			binding := reflect.TypeOf(model).Field(index).Tag.Get("binding")
			if strings.Contains(t, "omitempty") || !strings.Contains(binding, "required") {
				canOmit = true
			}
			t = strings.ReplaceAll(t, "omitempty", "")
			t = strings.ReplaceAll(t, ",", "")
			if canOmit {
				t += "?"
			}
			return t
		}
	}
	return ""
}
