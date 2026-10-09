package back

type op uint8

const (
	opAdd = iota
	opChg // change
	opDel // delete
)

type query struct {
	target string
	operator op
	sequence string
}

func opFromStr(in string) op {
	switch in {
	case "-": return opDel
	case "+": return opChg
	}
	return opAdd
}
