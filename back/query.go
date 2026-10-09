package back

import "regexp"

type op uint8

const (
	opAdd = iota
	opChg // change
	opDel // delete
)

var	isInsert =
regexp.MustCompile(`^\s*(\$\w+\s+=\s+"[а-яА-ЯёЁ\w]+(?:\s+[а-яА-ЯёЁ\w]+)*")|(\%\w+\s+=\s+\[\d+(,\d+)*\]\s+\:\s+"[а-яА-ЯёЁ\w]+(?:\s+[а-яА-ЯёЁ\w]+)*")\s*$`)

var isChange =
regexp.MustCompile(`^(\%\w+\s+\+\s+\$\w+)$`)

var isDelete =
regexp.MustCompile(`^(\%\w+\s+\-\s+\$\w+)$`)

type Query struct {
	target string
	operator op
	sequence string
}

func ValidateQuery(in string) bool {
	return isInsert.MatchString(in) || isChange.MatchString(in) || isDelete.MatchString(in)
}

func opFromStr(in string) op {
	switch in {
	case "-": return opDel
	case "+": return opChg
	}
	return opAdd
}

// func ParseQuery(in string) Query {
//
// }
