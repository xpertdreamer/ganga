package back

import (
	"errors"
	"ganga/util"
	"math/rand"
	"time"
)

type Answer int

const (
	AnswerUnknown Answer = iota
	AnswerYes
	AnswerNo					
)

type State struct {
	Id int64
	Percents map[string]float64
}

type Engine struct {
	Base *KBase
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
	return nil
}

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


func (e* Engine)RecalculatePercents(p Property, answer Answer) {
	if e.StatePoiner < uint64(len(e.States)-1) {
		e.States = e.States[:e.StatePoiner+1]
	}
	currentState := e.States[e.StatePoiner]
	nextPercents := make(map[string]float64, len(currentState.Percents))
	for k, v := range currentState.Percents {
		nextPercents[k] = v
	}
	var factor = 100.0 / float64(len(e.Base.Objects))
	var targetDivisor float64 = (100.0 - factor) / 100.0
	var totalWeight float64
	for objKey, currentPercent := range nextPercents {
		hasProp, err := e.Base.HasProperty(objKey, p.ID)
		if err != nil {
			util.Error("some error occured during recalculation")
			return
		}
		var multiplier float64 = 1.0
		switch answer {
		case AnswerYes:
			if hasProp {
				multiplier = 1.0 + (factor / 100.0)
			} else {
				multiplier = targetDivisor
			}
		case AnswerNo:
			if !hasProp {
				multiplier = 1.0 + (factor / 100.0)
			} else {
				multiplier = targetDivisor
			}
		}
		newPercent := currentPercent * multiplier
		nextPercents[objKey] = newPercent
		totalWeight += newPercent
	}

	if totalWeight > 0 {
		for objKey := range nextPercents {
			nextPercents[objKey] = (nextPercents[objKey] / totalWeight) * 100.0
		}
	}

	nextId := currentState.Id + 1

	nextState := State{
		Id:       nextId,
		Percents: nextPercents,
	}
	e.States = append(e.States, nextState)

	e.StatePoiner = uint64(nextId)
}

func (e* Engine)Start() error {
	if e == nil {
		return errors.New("cant deal with nil pointer engine")
	}
	initialNumQuest := len(e.Base.Properties)
	switch {
	case initialNumQuest > 50: initialNumQuest = int(float64(initialNumQuest) * 0.25)
	case initialNumQuest <= 50 && initialNumQuest > 5:	initialNumQuest = int(float64(initialNumQuest) * 0.5)
	}
	e.Questions = GetRandomKeys(e.Base, initialNumQuest)
	e.StatePoiner = 0
	util.Debug("initial questions: %d", initialNumQuest)
	count := len(e.Base.Objects)
	initialPercents := make(map[string]float64, count)
	var startValue float64 = 0.0
	if count > 0 {
		startValue = 100.0 / float64(count)
	}
	for objKey := range e.Base.Objects {
		initialPercents[objKey] = startValue
	}
	firstState := State{
		Id: 0,
		Percents: initialPercents,
	}
	e.States = append(e.States, firstState)
	return nil
}

func (e* Engine)GetResult() ([]string, error) {
	if e == nil {
		return nil, errors.New("cant deal with nil pointer engine")
	}
	currentPercents := e.States[e.StatePoiner].Percents

	count := len(e.Base.Objects)
	if count == 0 {
		return nil, errors.New("no objects in knowledge base")
	}

	startValue := 100.0 / float64(count)

	var resultIdx []string
	var max float64 = 0.0

	for _, val := range currentPercents {
		if val > max {
			max = val
		}
	}

	if max <= startValue {
		return nil, nil
	}

	for objKey, percent := range currentPercents {
		if int(percent) == int(max) {
			resultIdx = append(resultIdx, objKey)
		}
	}

	return resultIdx, nil
}

func (e* Engine)GetMoreQuestions(count int) int {
	if len(e.Base.Properties) == 0 || count <= 0 {
		return 0
	}

	existing := make(map[string]bool)
	for _, q := range e.Questions {
		existing[q] = true
	}

	var availableKeys []string
	for key := range e.Base.Properties {
		if !existing[key] {
			availableKeys = append(availableKeys, key)
		}
	}

	if len(availableKeys) == 0 {
		return 0
	}

	if count > len(availableKeys) {
		count = len(availableKeys)
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	addedCount := 0

	for i := 0; i < count; i++ {
		randIdx := r.Intn(len(availableKeys))
		e.Questions = append(e.Questions, availableKeys[randIdx])

		availableKeys[randIdx] = availableKeys[len(availableKeys)-1]
		availableKeys = availableKeys[:len(availableKeys)-1]
		addedCount++
	}

	return addedCount
}
