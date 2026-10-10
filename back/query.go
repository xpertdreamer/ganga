package back

import (
	"ganga/util"
	"regexp"
	"strings"
)

type op uint8

const (
	opAdd = iota
	opChg // change
	opDel // delete
	opDsf
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

type Query struct {
	target string
	operator op
	sequence string
}

func ValidateQuery(in string) bool {
	return isPropInsert.MatchString(in) || isObjInsert.MatchString(in) ||
		isChange.MatchString(in) || isDelete.MatchString(in) || isDeleteSelf.MatchString(in)
}

func ParseQuery(in string) Query {
	t := strings.TrimSpace(in)
	util.Debug("trimmed: %s", t)

	switch {
	case isDeleteSelf.MatchString(t):
		target := isDeleteSelf.FindStringSubmatch(t)
		return Query{
			target: strings.TrimRight(target[0], "-"),
			operator: opDsf,
			sequence: "",
		}

	case isChange.MatchString(t):
		tokens := strings.Fields(t)
		return Query{
			target: tokens[0],
			operator: opChg,
			sequence: tokens[2],
		}

	case isDelete.MatchString(t):
		tokens := strings.Fields(t)
		return Query {
			target: tokens[0],
			operator: opDel,
			sequence: tokens[2],
		}

	case isObjInsert.MatchString(t):
		parts := strings.SplitN(t, "=", 2)
		if len(parts) != 2 {
			return Query{}
		}
		return Query{
			target: strings.TrimSpace(parts[0]),
			operator: opAdd,
			sequence: strings.TrimSpace(parts[1]),
		}

	case isPropInsert.MatchString(t):
		parts := strings.SplitN(t, "=", 2)
		if len(parts) != 2 {
			return Query{}
		}
		return Query{
			target:   strings.TrimSpace(parts[0]),
			operator: opAdd,
			sequence: strings.TrimSpace(parts[1]),
		}
	}

	return Query{}
}
