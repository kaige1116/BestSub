package desc

import (
	"reflect"
)

type Data struct {
	Name    string `json:"name,omitempty"`
	Key     string `json:"key,omitempty"`
	Type    string `json:"type,omitempty"`
	Value   string `json:"value,omitempty"`
	Options string `json:"options,omitempty"`
	Require bool   `json:"require,omitempty"`
	Desc    string `json:"desc,omitempty"`
}

func Gen(v any) []Data {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return gen(t)
}

func gen(t reflect.Type) []Data {
	var items []Data
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Type.Kind() == reflect.Struct {
			items = append(items, gen(field.Type)...)
			continue
		}
		tag := field.Tag
		key, ok := tag.Lookup("json")
		if !ok {
			continue
		}
		typeName := tag.Get("type")
		if typeName == "" {
			typeName = getType(field.Type.Name())
		}
		item := Data{
			Name:    tag.Get("name"),
			Key:     key,
			Type:    typeName,
			Value:   tag.Get("value"),
			Options: tag.Get("options"),
			Require: tag.Get("require") == "true",
			Desc:    tag.Get("description"),
		}
		items = append(items, item)
	}
	return items
}

func getType(t string) string {
	switch t {
	case "bool":
		return "boolean"
	case "int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"float32", "float64":
		return "number"
	case "string", "[]byte":
		return "string"
	default:
		return "object"
	}
}
