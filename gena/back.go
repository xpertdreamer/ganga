package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
)

type object	struct {
	desc string
	tag string
	props map[string]bool
}

type model struct {
	objects map[string]*object
	props map[string]string
}

func newModel() *model {
	return &model {
		props: make(map[string]string),
		objects: make(map[string]*object),
	}
}

type rawProp struct {
	Description string `json:"description"`
}

type rawObj struct {
	Properties  []int  `json:"properties"`
	Description string `json:"description"`
	Tag string `json:"tag"`
}

type raw struct {
	Properties map[string]rawProp `json:"properties"`
	Objects map[string]rawObj `json:"objects"`
}

func (m* model) load(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("json:%w", err)
	}
	var raw raw
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("json:%w", err)
	}
	m.props = make(map[string]string, len(raw.Properties))
	for id, p := range raw.Properties {
		m.props[id] = p.Description
	}
	m.objects = make(map[string]*object, len(raw.Objects))
	for id, o := range raw.Objects {
		ob := &object {
			desc: o.Description,
			tag: o.Tag,
			props: make(map[string]bool, len(o.Properties)),
		}
		for _, n := range o.Properties {
			ob.props[strconv.FormatInt(int64(n), 10)] = true
		}
		m.objects[id] = ob
	}
	return nil
}

func (m* model) save(path string) error {
	raw := raw{
		Properties: make(map[string]rawProp, len(m.props)),
		Objects: make(map[string]rawObj, len(m.objects)),
	}
	for id, d := range m.props {
		raw.Properties[id] = rawProp{Description: d}
	}
	for id, o := range m.objects {
		props := make([]int, 0, len(o.props))
		for pid := range o.props {
			n, err := strconv.Atoi(pid)
			if err != nil {
				return fmt.Errorf("property %q from %q is not int", pid, id)
			}
			props = append(props, n)
		}
		sort.Ints(props)
		raw.Objects[id] = rawObj {
			Properties: props,
			Description: o.desc,
			Tag: o.tag,
		}
	}
	data, err := json.MarshalIndent(raw, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

func sortedKeys[V any](m map[string]V) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
