package render

import (
	"fmt"
	"path"
	"strings"

	"github.com/gdamore/tcell"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

// PersistentVolume renders a K8s PersistentVolume to screen.
type PersistentVolume struct{}

// ColorerFunc colors a resource row.
func (PersistentVolume) ColorerFunc() ColorerFunc {
	return func(ns string, r RowEvent) tcell.Color {
		c := DefaultColorer(ns, r)
		if r.Kind == EventAdd || r.Kind == EventUpdate {
			return c
		}

		status := strings.TrimSpace(r.Row.Fields[4])
		switch status {
		case "Bound":
			c = StdColor
		case "Available":
			c = tcell.ColorYellow
		default:
			c = ErrColor
		}

		return c
	}

}

// Header returns a header rbw.
func (PersistentVolume) Header(string) HeaderRow {
	return HeaderRow{
		Header{Name: "NAME"},
		Header{Name: "CAPACITY"},
		Header{Name: "ACCESS MODES"},
		Header{Name: "RECLAIM POLICY"},
		Header{Name: "STATUS"},
		Header{Name: "CLAIM"},
		Header{Name: "STORAGECLASS"},
		Header{Name: "REASON"},
		Header{Name: "AGE", Decorator: AgeDecorator},
	}
}

// Render renders a K8s resource to screen.
func (p PersistentVolume) Render(o interface{}, ns string, r *Row) error {
	raw, ok := o.(*unstructured.Unstructured)
	if !ok {
		return fmt.Errorf("Expected PersistentVolume, but got %T", o)
	}
	var pv v1.PersistentVolume
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(raw.Object, &pv)
	if err != nil {
		return err
	}

	phase := pv.Status.Phase
	if pv.ObjectMeta.DeletionTimestamp != nil {
		phase = "Terminating"
	}
	var claim string
	if pv.Spec.ClaimRef != nil {
		claim = path.Join(pv.Spec.ClaimRef.Namespace, pv.Spec.ClaimRef.Name)
	}
	class, found := pv.Annotations[v1.BetaStorageClassAnnotation]
	if !found {
		class = pv.Spec.StorageClassName
	}

	size := pv.Spec.Capacity[v1.ResourceStorage]

	r.ID = MetaFQN(pv.ObjectMeta)
	r.Fields = Fields{
		pv.Name,
		size.String(),
		accessMode(pv.Spec.AccessModes),
		string(pv.Spec.PersistentVolumeReclaimPolicy),
		string(phase),
		claim,
		class,
		pv.Status.Reason,
		toAge(pv.ObjectMeta.CreationTimestamp),
	}

	return nil
}

// ----------------------------------------------------------------------------
// Helpers...

func accessMode(aa []v1.PersistentVolumeAccessMode) string {
	dd := accessDedup(aa)
	s := make([]string, 0, len(dd))
	for i := 0; i < len(aa); i++ {
		switch {
		case accessContains(dd, v1.ReadWriteOnce):
			s = append(s, "RWO")
		case accessContains(dd, v1.ReadOnlyMany):
			s = append(s, "ROX")
		case accessContains(dd, v1.ReadWriteMany):
			s = append(s, "RWX")
		}
	}

	return strings.Join(s, ",")
}

func accessContains(cc []v1.PersistentVolumeAccessMode, a v1.PersistentVolumeAccessMode) bool {
	for _, c := range cc {
		if c == a {
			return true
		}
	}

	return false
}

func accessDedup(cc []v1.PersistentVolumeAccessMode) []v1.PersistentVolumeAccessMode {
	set := []v1.PersistentVolumeAccessMode{}
	for _, c := range cc {
		if !accessContains(set, c) {
			set = append(set, c)
		}
	}

	return set
}
