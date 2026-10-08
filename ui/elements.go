package ui

import (
	"fmt"
	"ganga/back"
	"ganga/util"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

func CreatePreview(k* back.KBase) string {
	var sb strings.Builder
	sb.Grow(512)

	sb.WriteString("Objects : \n\n")
	for _, obj := range k.Objects {
		sb.WriteString("\t")
		sb.WriteString(obj.Description)
		sb.WriteString(" ")
		fmt.Fprintf(&sb, "%v", obj.Properties)
		sb.WriteString(",\n\n")
	}

	sb.WriteString("Properties : \n\n")
	for _, prop := range k.Properties {
		sb.WriteString("\t")
		sb.WriteString(prop.ID)
		sb.WriteString(":")
		sb.WriteString(prop.Description)
		sb.WriteString(",\n")
	}

	return sb.String()
}

type MascotPaths struct {
	Reflective fyne.Resource
	Happy fyne.Resource
	Stable fyne.Resource
}

func LoadMaskotCache() (MascotPaths, error) {
	var cache MascotPaths
	exePath, err := os.Executable()
	if err != nil {
		return cache, fmt.Errorf("cant recognize path to binary: %w", err)
	}

	buildDir := filepath.Dir(exePath)
	assetsDir := filepath.Join(buildDir, "..", "assets")

	files := []struct {
		name     string
		fileName string
		target   *fyne.Resource
	}{
		{"reflective", "maskot_refl.png", &cache.Reflective},
		{"happy", "maskot_happy.png", &cache.Happy},
		{"stable", "maskot_stab.png", &cache.Stable},
	}

	for _, file := range files {
		fullPath := filepath.Join(assetsDir, file.fileName)

		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			return cache, fmt.Errorf("asset file %s is missing in: %s", file.name, fullPath)
		}

		data, err := os.ReadFile(fullPath)
		if err != nil {
			return cache, fmt.Errorf("cant read file %s: %w", file.name, err)
		}

		*file.target = fyne.NewStaticResource(file.fileName, data)
	}

	return cache, nil
}

func ChooseMaskotRand(maskot *canvas.Image, cache MascotPaths) {
	if maskot == nil {
		return
	}
	availableResources := []fyne.Resource{cache.Reflective, cache.Happy, cache.Stable}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomIndex := r.Intn(len(availableResources))
	maskot.Resource = availableResources[randomIndex]
	maskot.Refresh()
}

func HandleButtonNext(e* back.Engine, quest *widget.Label, numQuest *widget.Label, answer back.Answer) bool {
	if e == nil {
		util.Error("cant deal with nil pointers")
		return false
	}
	util.Debug("state pointer: %d", e.StatePoiner)
	util.Debug("questions: %v", e.Questions)
	if e.StatePoiner >= uint64(len(e.Questions)) {
		util.Debug("All questions are already answered")
		return false
	}

	currentKey := e.Questions[e.StatePoiner]
	currentProp := e.Base.Properties[currentKey]

	e.RecalculatePercents(currentProp, answer)

	if e.StatePoiner < uint64(len(e.Questions)) {
		nextKey := e.Questions[e.StatePoiner]
		nextProp := e.Base.Properties[nextKey]
		quest.SetText(nextProp.Description + "?")
		numQuest.SetText(strconv.FormatUint(e.StatePoiner + 1, 10))
	} else {
		quest.SetText("")
		numQuest.SetText("")
		return true
	}
	util.Debug("new state pointer: %d", e.StatePoiner)
	util.Debug("recalculated percents: %v", e.States[e.StatePoiner].Percents)

	return false
}

func HandleButtonPrev(e* back.Engine, quest *widget.Label, numQuest *widget.Label) {
	if e == nil {
		return
	}
	if e.StatePoiner == 0 {
		util.Debug("Already at the first question")
		return
	}
	e.StatePoiner -= 1
	prevKey := e.Questions[e.StatePoiner]
	prevProp := e.Base.Properties[prevKey]
	quest.SetText(prevProp.Description + "?")
	numQuest.SetText(strconv.FormatUint(e.StatePoiner + 1, 10))
	util.Debug("prev pressed. state pointer: %d, percents: %v", e.StatePoiner, e.States[e.StatePoiner].Percents)
}
