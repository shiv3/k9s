package dao

import (
	"bytes"
	"errors"
	"math"

	"github.com/derailed/tview"
	runewidth "github.com/mattn/go-runewidth"
	"github.com/rs/zerolog/log"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/cli-runtime/pkg/printers"
)

func toPerc(v1, v2 float64) float64 {
	if v2 == 0 {
		return 0
	}
	return math.Round((v1 / v2) * 100)
}

// Truncate a string to the given l and suffix ellipsis if needed.
func Truncate(str string, width int) string {
	return runewidth.Truncate(str, width, string(tview.SemigraphicsHorizontalEllipsis))
}

// ToYAML converts a resource to its YAML representation.
func ToYAML(o runtime.Object) (string, error) {
	if o == nil {
		return "", errors.New("no object to yamlize")
	}
	var (
		buff bytes.Buffer
		p    printers.YAMLPrinter
	)
	err := p.PrintObj(o, &buff)
	if err != nil {
		log.Error().Msgf("Marshal Error %v", err)
		return "", err
	}

	return buff.String(), nil
}
