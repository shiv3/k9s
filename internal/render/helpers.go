package render

import (
	"path"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/derailed/k9s/internal/client"
	"github.com/derailed/tview"
	runewidth "github.com/mattn/go-runewidth"
	"github.com/rs/zerolog/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/duration"
)

const megaByte = 1024 * 1024

// ToMB converts bytes to megabytes.
func ToMB(v int64) float64 {
	return float64(v) / megaByte
}

func asSelector(s *metav1.LabelSelector) string {
	sel, err := metav1.LabelSelectorAsSelector(s)
	if err != nil {
		log.Error().Err(err).Msg("Selector conversion failed")
		return NAValue
	}

	return sel.String()
}

type metric struct {
	cpu, mem string
}

func noMetric() metric {
	return metric{cpu: NAValue, mem: NAValue}
}

// MetaFQN returns a fully qualified resource name.
func MetaFQN(m metav1.ObjectMeta) string {
	if m.Namespace == "" {
		return FQN(client.ClusterScope, m.Name)
	}

	return FQN(m.Namespace, m.Name)
}

// FQN returns a fully qualified resource name.
func FQN(ns, n string) string {
	if ns == "" {
		return n
	}
	return ns + "/" + n
}

// ToSelector flattens a map selector to a string selector.
func toSelector(m map[string]string) string {
	s := make([]string, 0, len(m))
	for k, v := range m {
		s = append(s, k+"="+v)
	}

	return strings.Join(s, ",")
}

// Blank checks if a collection is empty or all values are blank.
func blank(s []string) bool {
	for _, v := range s {
		if len(v) != 0 {
			return false
		}
	}
	return true
}

// Join a slice of strings, skipping blanks.
func join(a []string, sep string) string {
	switch len(a) {
	case 0:
		return ""
	case 1:
		return a[0]
	}

	var b []string
	for _, s := range a {
		if s != "" {
			b = append(b, s)
		}
	}
	if len(b) == 0 {
		return ""
	}

	n := len(sep) * (len(b) - 1)
	for i := 0; i < len(b); i++ {
		n += len(a[i])
	}

	var buff strings.Builder
	buff.Grow(n)
	buff.WriteString(a[0])
	for _, s := range b[1:] {
		buff.WriteString(sep)
		buff.WriteString(s)
	}

	return buff.String()
}

// AsPerc prints a number as a percentage.
func AsPerc(f float64) string {
	return strconv.Itoa(int(f))
}

// ToPerc computes the ratio of two numbers as a percentage.
func toPerc(v1, v2 float64) float64 {
	if v2 == 0 {
		return 0
	}
	return (v1 / v2) * 100
}

// Namespaced return a namesapace and a name.
func Namespaced(n string) (string, string) {
	ns, po := path.Split(n)

	return strings.Trim(ns, "/"), po
}

func missing(s string) string {
	return check(s, MissingValue)
}

func na(s string) string {
	return check(s, NAValue)
}

func check(s, sub string) string {
	if len(s) == 0 {
		return sub
	}

	return s
}

func boolToStr(b bool) string {
	switch b {
	case true:
		return "true"
	default:
		return "false"
	}
}

func toAge(timestamp metav1.Time) string {
	return time.Since(timestamp.Time).String()
}

func toAgeHuman(s string) string {
	d, err := time.ParseDuration(s)
	if err != nil {
		return NAValue
	}

	return duration.HumanDuration(d)
}

// Truncate a string to the given l and suffix ellipsis if needed.
func Truncate(str string, width int) string {
	return runewidth.Truncate(str, width, string(tview.SemigraphicsHorizontalEllipsis))
}

func mapToStr(m map[string]string) (s string) {
	if len(m) == 0 {
		return MissingValue
	}

	kk := make([]string, 0, len(m))
	for k := range m {
		kk = append(kk, k)
	}
	sort.Strings(kk)

	for i, k := range kk {
		s += k + "=" + m[k]
		if i < len(kk)-1 {
			s += ","
		}
	}

	return
}

// ToMillicore shows cpu reading for human.
func ToMillicore(v int64) string {
	return strconv.Itoa(int(v))
}

// ToMi shows mem reading for human.
func ToMi(v float64) string {
	return strconv.Itoa(int(v))
}

func boolPtrToStr(b *bool) string {
	if b == nil {
		return "false"
	}

	return boolToStr(*b)
}

// Check if string is in a string list.
func in(ll []string, s string) bool {
	for _, l := range ll {
		if l == s {
			return true
		}
	}
	return false
}

// Pad a string up to the given length or truncates if greater than length.
func Pad(s string, width int) string {
	if len(s) == width {
		return s
	}

	if len(s) > width {
		return Truncate(s, width)
	}

	return s + strings.Repeat(" ", width-len(s))
}
