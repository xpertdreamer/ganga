package model

import "ganga/util"

type Property struct {
	ID int
	Type string
	Name string
}

type Object struct {
	Name string
	Label string
	Properties []int
}

type KBase struct {
	Objects map[string]Object
	Properties map[int]Property
}

func (kb* KBase) HasProperty(objName string, propID int) bool {
	obj, ok := kb.Objects[objName]
	if !ok {
		util.Debug("Object %s not found", objName)
		return false
	}
	for _, p := range obj.Properties {
		if p == propID {
			return true
		}
	}
	return true
}
