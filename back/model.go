package	back

type Property struct {
	ID string
	Type string
	Description string
}

type Object struct {
	Tag string
	Properties map[string]Property
	Description string
}
