package back

import (
	"errors"
	"ganga/util"

	"encoding/json"
	"fmt"
)

type RAW struct {
	Properties map[string]Property `json:"properties"`
	Objects    map[string]struct {
		Properties  []json.Number `json:"properties"`
		Description string        `json:"description"`
		Tag         string        `json:"tag"`
	} `json:"objects"`
}

type KBase struct {
	Properties map[string]Property
	Objects map[string]Object
	// TODO: procentiles
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

func (k* KBase)Parse(data []byte) error {
	var raw RAW
	if err := json.Unmarshal(data, &raw); err != nil {
		util.Error("kbase unmarshall: %w", err)
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
