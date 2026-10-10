package back

import (
	"fmt"
	"ganga/util"
	"regexp"
	"strings"
)

type op uint8

const (
	opAdd = iota
	opChg // change
	opDel // delete
	opDsf // delete self
)

var	isPropInsert =
regexp.MustCompile(`^\s*\$\w+\s+=\s+"[а-яА-ЯёЁ\w]+(?:\s+[а-яА-ЯёЁ\w]+)*"\s*$`)

var isObjInsert =
regexp.MustCompile(`^\s*\%\w+\s+=\s+\[\d+(,\d+)*\]\s+\:\s+"[а-яА-ЯёЁ\w]+(?:\s+[а-яА-ЯёЁ\w]+)*"\s*$`)

var isChange =
regexp.MustCompile(`^(\%\w+\s+\+\s+\$\w+)$`)

var isDelete =
regexp.MustCompile(`^(\%\w+\s+\-\s+\$\w+)$`)

var isDeleteSelf =
regexp.MustCompile(`^([\$\%]\w+\-\s*)$`)

var digit =
regexp.MustCompile(`\d+`)

type QueryRaw struct {
	target string
	operator op
	sequence string
}

func ValidateQuery(in string) bool {
	return isPropInsert.MatchString(in) || isObjInsert.MatchString(in) ||
		isChange.MatchString(in) || isDelete.MatchString(in) || isDeleteSelf.MatchString(in)
}

func ParseQuery(in string) QueryRaw {
	t := strings.TrimSpace(in)
	util.Debug("trimmed: %s", t)

	switch {
	case isDeleteSelf.MatchString(t):
		target := isDeleteSelf.FindStringSubmatch(t)
		return QueryRaw{
			target: strings.TrimRight(target[0], "-"),
			operator: opDsf,
			sequence: "",
		}

	case isChange.MatchString(t):
		tokens := strings.Fields(t)
		return QueryRaw{
			target: tokens[0],
			operator: opChg,
			sequence: tokens[2],
		}

	case isDelete.MatchString(t):
		tokens := strings.Fields(t)
		return QueryRaw {
			target: tokens[0],
			operator: opDel,
			sequence: tokens[2],
		}

	case isObjInsert.MatchString(t):
		parts := strings.SplitN(t, "=", 2)
		if len(parts) != 2 {
			return QueryRaw{}
		}
		return QueryRaw{
			target: strings.TrimSpace(parts[0]),
			operator: opAdd,
			sequence: strings.TrimSpace(parts[1]),
		}

	case isPropInsert.MatchString(t):
		parts := strings.SplitN(t, "=", 2)
		if len(parts) != 2 {
			return QueryRaw{}
		}
		return QueryRaw{
			target:   strings.TrimSpace(parts[0]),
			operator: opAdd,
			sequence: strings.TrimSpace(parts[1]),
		}
	}

	return QueryRaw{}
}

func Submit(e *Engine, q QueryRaw) error {
	if e == nil {
		return fmt.Errorf("cant deal with engine nil pointer")
	}
	switch q.operator {
	case opAdd:
		return submitAdd(e, q)
	case opChg:
		return submitChg(e, q)
	case opDel:
		return submitDel(e, q)
	case opDsf:
		return submitDsf(e, q)
	}
	return fmt.Errorf("unknown operator")
}

func hasAllKeys[V any](m map[string]V, keys []string) bool {
	for _, key := range keys {
		if _, exists := m[key]; !exists {
			return false
		}
	}
	return true
}

func submitAdd(e *Engine, q QueryRaw) error {
	switch {
	case strings.HasPrefix(q.target, "$"):
		target := strings.TrimPrefix(q.target, "$")
		second := strings.TrimPrefix(q.sequence, "\"")
		second = strings.TrimSuffix(second, "\"")
		newProp := Property {
			ID: target,
			Type: "bool",
			Description: second,
		}
		e.Base.Properties[target] = newProp
		return nil

	case strings.HasPrefix(q.target, "%"):
		target := strings.TrimPrefix(q.target, "%")
		parts := strings.Split(q.sequence, " : ")
		if len(parts) != 2 {
			return fmt.Errorf("bad sequence %s", q.sequence)
		}
		left := parts[0]
		right := parts[1]
		digits := digit.FindAllString(left, -1)

		if !hasAllKeys(e.Base.Properties, digits) {
			return fmt.Errorf("bad sequence %s", q.sequence)
		}

		right = strings.TrimPrefix(right, "\"")
		right = strings.TrimSuffix(right, "\"")

		props := make(map[string]Property)
		for _, d := range digits {
			if prop, exists := e.Base.Properties[d]; exists {
				props[d] = prop
			}
		}

		newObj := Object {
			Tag: target,
			Properties: props,
			Description: right,
		}

		e.Base.Objects[target] = newObj

		return nil
	}
	return fmt.Errorf("invalid query %v", q)
}

func submitChg(e *Engine, q QueryRaw) error {
	target := strings.TrimPrefix(q.target, "%")
	second := strings.TrimPrefix(q.sequence, "$")
	obj, ok := e.Base.Objects[target]
	if !ok {
		return fmt.Errorf("object %q not found", q.target)
	}
	var prop *Property
	for _, p := range e.Base.Properties {
		if p.ID == second {
			prop = &p
			break
		}
	}
	if prop == nil {
		return fmt.Errorf("property %q not found", q.sequence)
	}
	obj.Properties[second] = *prop
	return nil
}

func submitDel(e *Engine, q QueryRaw) error {
	target := strings.TrimPrefix(q.target, "%")
	second := strings.TrimPrefix(q.sequence, "$")
	obj, ok := e.Base.Objects[target]
	if !ok {
		return fmt.Errorf("object %q not found", q.target)
	}

	if _, ok := obj.Properties[second]; !ok {
		return fmt.Errorf("property %q not present in %q", q.sequence, q.target)
	}
	delete(obj.Properties, second)
	return nil
}

func submitDsf(e *Engine, q QueryRaw) error {
	switch {
	case strings.HasPrefix(q.target, "%"):
		name := strings.TrimPrefix(q.target, "%")
		if _, ok := e.Base.Objects[name]; !ok {
			return fmt.Errorf("object %q not found", q.target)
		}
		delete(e.Base.Objects, name)
		return nil

	case strings.HasPrefix(q.target, "$"):
		name := strings.TrimPrefix(q.target, "$")
		var removed bool = false
		for _, obj := range e.Base.Objects {
			if obj.Properties == nil {
				continue
			}
			if _, ok := obj.Properties[name]; ok {
				delete(obj.Properties, name)
				removed = true
			}
		}
		if !removed {
			return fmt.Errorf("property %q not found", name)
		}
		delete(e.Base.Properties, name)
		return nil
	}

	return fmt.Errorf("invalid target %q", q.target)
}
