package back

import (
	"errors"
	"ganga/util"

	"encoding/json"
	"fmt"
	"sort"
	"io"
)

type RawObject struct {
	Properties  []json.Number `json:"properties"`
	Description string        `json:"description"`
	Tag         string        `json:"tag"`
}

type RAW struct {
	Properties map[string]Property  `json:"properties"`
	Objects    map[string]RawObject `json:"objects"`
}

type KBase struct {
	Properties map[string]Property
	Objects map[string]Object
}

func (k *KBase) toRAW() RAW {
	raw := RAW{
		Properties: make(map[string]Property, len(k.Properties)),
		Objects:    make(map[string]RawObject, len(k.Objects)),
	}

	for id, p := range k.Properties {
		raw.Properties[id] = p
	}

	for id, obj := range k.Objects {
		props := make([]json.Number, 0, len(obj.Properties))
		for pid := range obj.Properties {
			props = append(props, json.Number(pid))
		}

		sort.Slice(props, func(i, j int) bool {
			return props[i] < props[j]
		})

		raw.Objects[id] = RawObject{
			Tag:         obj.Tag,
			Description: obj.Description,
			Properties:  props,
		}
	}

	return raw
}

func (k *KBase) Save(w io.Writer) error {
	if k == nil {
		return errors.New("kbase is nil")
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(k.toRAW())
}

func (k* KBase)HasProperty(objId string, propId string) (bool, error) {
	if obj, ok := k.Objects[objId]; ok {
		if _, found := obj.Properties[propId]; found {
			util.Debug("(HasProperty) found [%s] for [%s]", objId, propId)
			return true, nil
		}
		return false, nil
	} else {
		util.Error("(HasProperty) Object [%s] not found!", objId)
		return false, errors.New("object not found")
	}
}

func(k* KBase)Clear() {
	if k != nil {
		clear(k.Objects)
		clear(k.Properties)
	}
}

func (k* KBase)Parse(data []byte) error {
	var raw RAW
	if err := json.Unmarshal(data, &raw); err != nil {
		util.Error("kbase unmarshall: %s", err.Error())
		return errors.New("kbase unmarshall error")
	}
	k.Properties = make(map[string]Property, len(raw.Properties))
	for id, p := range raw.Properties {
		p.ID = id
		k.Properties[id] = p
	}
	k.Objects = make(map[string]Object, len(raw.Objects))
	for objID, ro := range raw.Objects {
		obj := Object{
			Tag:         ro.Tag,
			Description: ro.Description,
			Properties:  make(map[string]Property, len(ro.Properties)),
		}
		for _, n := range ro.Properties {
			pid := n.String()
			p, ok := k.Properties[pid]
			if !ok {
				util.Error("object %q references unknown property %q", objID, pid)
				return fmt.Errorf("object %q references unknown property %q", objID, pid)
			}
			obj.Properties[pid] = p
		}
		k.Objects[objID] = obj
		util.Debug("(Parse) loaded object [%s] with %d properties", objID, len(obj.Properties))
	}
	return nil
}

func NewKBase() *KBase {
	return &KBase{
		Properties: make(map[string]Property),
		Objects:    make(map[string]Object),
	}
}
