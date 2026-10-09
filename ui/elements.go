package ui

import (
	"fmt"
	"ganga/back"
	"ganga/util"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const (
	MascotCatRefl = "maskot_cat_refl.png"
	MascotCatHappy = "maskot_cat_happy.png"
	MascotCatStab = "maskot_cat_stab.png"

	MascotGopherRefl = "maskot_refl.png"
	MascotGopherHappy = "maskot_happy.png"
	MascotGopherStab = "maskot_stab.png"
)

type Preset struct {
	name     string
	fileName string
	target   *fyne.Resource
}

func newTable(n int, headers []string, row func(i int) []string) *widget.Table {
	t := widget.NewTable(
		func() (int, int) { return n + 1, len(headers) },
		func() fyne.CanvasObject {
			l := widget.NewLabel("")
			l.Wrapping = fyne.TextWrapWord
			return l
		},
		func(id widget.TableCellID, o fyne.CanvasObject) {
			l := o.(*widget.Label)
			if id.Row == 0 {
				l.TextStyle = fyne.TextStyle{Bold: true}
				l.SetText(headers[id.Col])
				return
			}
			l.TextStyle = fyne.TextStyle{}
			l.SetText(row(id.Row - 1)[id.Col])
		},
	)

	t.SetColumnWidth(0, 200)
	t.SetColumnWidth(1, 400)
	if len(headers) > 2 {
		t.SetColumnWidth(2, 400)
	}
	t.Resize(fyne.NewSize(900, 500))
	return t
}

func propKeysOfObject(o back.Object) string {
	keys := make([]string, 0)
	for k := range o.Properties {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return strings.Join(keys, ", ")
}

func CreatePreview(k* back.KBase) fyne.CanvasObject {
	propKeys := make([]string, 0, len(k.Properties))
	for i := range k.Properties {
		propKeys = append(propKeys, i)
	}
	sort.Strings(propKeys)

	objKeys := make([]string, 0, len(k.Objects))
	for i := range k.Objects {
		objKeys = append(objKeys, i)
	}
	sort.Strings(objKeys)

	propTable := newTable(
		len(propKeys),
		[]string {"Ключ", "Описание"},
		func(i int) []string {
			j := propKeys[i]
			return []string{j, k.Properties[j].Description}
		},
	)
	objTable := newTable(
		len(objKeys),
		[]string {"Ключ", "Описание", "Свойства"},
		func(i int) []string {
			j := objKeys[i]
			return []string {j, k.Objects[j].Description, propKeysOfObject(k.Objects[j])}
		},
	)

	nav := container.NewAppTabs(
		container.NewTabItem("Свойства ("+strconv.Itoa(len(propKeys))+" )", propTable),
		container.NewTabItem("Объекты ("+strconv.Itoa(len(objKeys))+" )", objTable),
	)

	return nav
}

type MascotPaths struct {
	Reflective fyne.Resource
	Happy fyne.Resource
	Stable fyne.Resource
}

func LoadMaskotCache(preset uint8) (MascotPaths, error) {
	var cache MascotPaths
	exePath, err := os.Executable()
	if err != nil {
		return cache, fmt.Errorf("cant recognize path to binary: %w", err)
	}

	buildDir := filepath.Dir(exePath)
	assetsDir := filepath.Join(buildDir, "..", "assets")

	var pres []Preset

	var PresetGopher = []Preset {
		{"reflective", MascotGopherRefl, &cache.Reflective},
			{"happy", MascotGopherHappy, &cache.Happy},
			{"stable", MascotGopherStab, &cache.Stable},
		}

	var PresetCat = []Preset {
		{"reflective", MascotCatRefl, &cache.Reflective},
			{"happy", MascotCatHappy, &cache.Happy},
			{"stable", MascotCatStab, &cache.Stable},
		}

	switch preset {
	case 0: pres = PresetGopher
	case 1: pres = PresetCat
	}
	files := pres

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
