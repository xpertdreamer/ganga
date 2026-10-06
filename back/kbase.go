package back

import (
	"errors"
	"ganga/util"
)

type KBase struct {
	Properties map[string]Property
	Objects map[string]Object
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
