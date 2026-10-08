package back

import (
	"errors"
	"ganga/util"
	"math/rand"
	"time"
)

type State struct {
	Id int64
	Percents map[string]float64
}

type Engine struct {
	Base *KBase
	run bool
	Questions []string
	States []State
	StatePoiner uint64
}

func (e* Engine)Create(k* KBase) error {
	if e.Base != nil {
		return errors.New("you already initialized engine")
	}
	if k == nil {
		return errors.New("cant create engine from nil pointer")
	}
	e.Base = k
	e.run = false
	return nil
}

// TODO: in func { traverse 'Percents' and check if Object in KBase with key [string] have property 'x' -> do smth with procentiles}

func GetRandomKeys(kb *KBase, count int) []string {
	if len(kb.Properties) == 0 || count <= 0 {
		return nil
	}
	allKeys := make([]string, 0, len(kb.Properties))
	for key := range kb.Properties {
		allKeys = append(allKeys, key)
	}
	if count > len(allKeys) {
		count = len(allKeys)
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	result := make([]string, 0, count)
	for i := 0; i < count; i++ {
		randIdx := r.Intn(len(allKeys))
		result = append(result, allKeys[randIdx])
		allKeys[randIdx] = allKeys[len(allKeys)-1]
		allKeys = allKeys[:len(allKeys)-1]
	}
	return result
}


// func (e* Engine)NextQuestion()

func (e* Engine)Start() error {
	if e == nil {
		return errors.New("cant deal with nil pointer engine")
	}
	e.run = true
	initialNumQuest := len(e.Base.Properties)
	switch {
	case initialNumQuest > 50: initialNumQuest = int(float64(initialNumQuest) * 0.25)
	case initialNumQuest <= 50 && initialNumQuest > 5:	initialNumQuest = int(float64(initialNumQuest) * 0.5)
	}
	e.Questions = GetRandomKeys(e.Base, initialNumQuest)
	e.StatePoiner = 0
	util.Debug("initial questions: %d", initialNumQuest)
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
