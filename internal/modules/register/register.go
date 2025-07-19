package register

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"strings"

	"github.com/bestruirui/bestsub/internal/utils/log"
)

type instance interface {
	Init() error
	Exec(ctx context.Context, log *log.Logger, args ...any) error
}

type desc struct {
	Name    string `json:"name,omitempty"`
	Type    string `json:"type,omitempty"`
	Default string `json:"default,omitempty"`
	Options string `json:"options,omitempty"`
	Require bool   `json:"require,omitempty"`
	Desc    string `json:"desc,omitempty"`
}

type registerInfo struct {
	im  map[string]instance
	aim map[string][]desc
}

var registers = map[string]*registerInfo{}

func register(t string, m string, i instance) {
	r, exists := registers[t]
	if !exists {
		r = &registerInfo{
			im:  make(map[string]instance),
			aim: make(map[string][]desc),
		}
		registers[t] = r
	}

	r.im[m] = i
	r.aim[m] = genDesc(reflect.TypeOf(i).Elem())
}

func get(t string, m string, c string) (instance, error) {
	ri, exists := registers[t]
	if !exists {
		return nil, errors.New("category not found")
	}

	info, exists := ri.im[m]
	if !exists {
		return nil, errors.New("item not found")
	}

	if c != "" {
		err := json.Unmarshal([]byte(c), info)
		if err != nil {
			return nil, err
		}
	}

	return info, nil
}

func getList(t string) []string {
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

func getInfoMap(t string) map[string][]desc {
	ri, exists := registers[t]
	if !exists {
		return nil
	}
	return ri.aim
}

func genDesc(t reflect.Type) []desc {
	var items []desc
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
		item := desc{
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
