package register

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"
)

type Desc struct {
	Name    string `json:"name,omitempty"`
	Type    string `json:"type,omitempty"`
	Default string `json:"default,omitempty"`
	Options string `json:"options,omitempty"`
	Require bool   `json:"require,omitempty"`
	Desc    string `json:"desc,omitempty"`
}

type registerInfo struct {
	im  map[string]any
	aim map[string][]Desc
}

var registers = map[string]*registerInfo{}

func register(t string, m string, i any) {
	r, exists := registers[t]
	if !exists {
		r = &registerInfo{
			im:  make(map[string]any),
			aim: make(map[string][]Desc),
		}
		registers[t] = r
	}

	r.im[m] = i
	r.aim[m] = genDesc(reflect.TypeOf(i).Elem())
}

func Get[T any](t string, m string, c string) (T, error) {
	ri, exists := registers[t]
	if !exists {
		return *new(T), errors.New("category not found")
	}

	info, exists := ri.im[m]
	if !exists {
		return *new(T), errors.New("item not found")
	}

	if c != "" {
		err := json.Unmarshal([]byte(c), info)
		if err != nil {
			return *new(T), err
		}
	}

	return info.(T), nil
}

func GetList(t string) []string {
	ri, exists := registers[t]
	if !exists {
		return nil
	}

	keys := make([]string, 0, len(ri.im))
	for k := range ri.im {
		keys = append(keys, k)
	}
	return keys
}

func GetInfoMap(t string) map[string][]Desc {
	ri, exists := registers[t]
	if !exists {
		return nil
	}
	return ri.aim
}

func genDesc(t reflect.Type) []Desc {
	var items []Desc
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Type.Kind() == reflect.Struct {
			items = append(items, genDesc(field.Type)...)
			continue
		}
		tag := field.Tag
		name, ok := tag.Lookup("desc")
		if !ok {
			continue
		}
		item := Desc{
			Name:    name,
			Type:    strings.ToLower(field.Type.Name()),
			Default: tag.Get("default"),
			Options: tag.Get("options"),
			Require: tag.Get("require") == "true",
			Desc:    tag.Get("description"),
		}
		items = append(items, item)
	}
	return items
}
