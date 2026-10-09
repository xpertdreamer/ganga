package main

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
