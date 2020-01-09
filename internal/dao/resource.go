package dao

import (
	"context"
	"fmt"

	"github.com/derailed/k9s/internal"
	"github.com/rs/zerolog/log"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
)

var (
	_ Accessor  = (*Resource)(nil)
	_ Describer = (*Resource)(nil)
	_ Nuker     = (*Resource)(nil)
)

// Resource represents an informer based resource.
type Resource struct {
	Generic
}

// List returns a collection of resources.
func (r *Resource) List(ctx context.Context, ns string) ([]runtime.Object, error) {
	log.Debug().Msgf("INF-LIST %q:%q", ns, r.gvr)
	strLabel, ok := ctx.Value(internal.KeyLabels).(string)
	lsel := labels.Everything()
	if sel, err := labels.ConvertSelectorToLabelsMap(strLabel); ok && err == nil {
		lsel = sel.AsSelector()
	}

	return r.Factory.List(r.gvr.String(), ns, false, lsel)
}

// Get returns a resource instance if found, else an error.
func (r *Resource) Get(ctx context.Context, path string) (runtime.Object, error) {
	return r.Factory.Get(r.gvr.String(), path, true, labels.Everything())
}

// ToYAML returns a resource yaml.
func (r *Resource) ToYAML(path string) (string, error) {
	o, err := r.Get(context.Background(), path)
	if err != nil {
		return "", err
	}

	raw, err := ToYAML(o)
	if err != nil {
		return "", fmt.Errorf("unable to marshal resource %s", err)
	}
	return raw, nil
}
