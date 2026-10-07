package back

import (
	"errors"
)

type Engine struct {
	base *KBase
	run bool
}

func (e* Engine)Create(k* KBase) error {
	if e.base != nil {
		return errors.New("you already initialized engine")
	}
	if k == nil {
		return errors.New("cant create engine from nil pointer")
	}
	e.base = k
	e.run = false
	return nil
}

func (e* Engine)Start() error {
	if e == nil {
		return errors.New("cant deal with nil pointer engine")
	}
	e.run = true
	return nil
}

func (e* Engine)Stop() error {
	if e == nil {
		return errors.New("cant deal with nil pointer engine")
	}
	e.run = false
	return nil
}

func (e* Engine)IsRunning() bool {
	return e.run
}
