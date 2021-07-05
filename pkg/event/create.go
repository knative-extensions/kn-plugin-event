package event

import (
	"errors"
	"fmt"
	"strings"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/wavesoftware/go-ensure"
)

var (
	// ErrCantMarshalAsJSON is returned if given CE data can't be marshalled
	// as JSON.
	ErrCantMarshalAsJSON = errors.New("can't marshal as JSON")

	// ErrCantSetField is returned if given field can't be applied.
	ErrCantSetField = errors.New("can't set field")
)

// NewDefault creates a default CloudEvent.
func NewDefault() *cloudevents.Event {
	e := cloudevents.NewEvent()
	e.SetType(DefaultType)
	e.SetID(NewID())
	ensure.NoError(e.SetData(cloudevents.ApplicationJSON, map[string]string{}))
	e.SetSource(DefaultSource())
	e.SetTime(time.Now())
	ensure.NoError(e.Validate())
	return &e
}

// CreateFromSpec will create an event by parsing given args.
func CreateFromSpec(spec *Spec) (*cloudevents.Event, error) {
	e := NewDefault()
	e.SetID(spec.ID)
	e.SetSource(spec.Source)
	e.SetType(spec.Type)
	m := map[string]interface{}{}
	for _, fieldSpec := range spec.Fields {
		if err := updateMapWithSpec(m, fieldSpec); err != nil {
			return nil, err
		}
	}
	err := e.SetData(cloudevents.ApplicationJSON, m)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCantMarshalAsJSON, err)
	}
	return e, nil
}

func updateMapWithSpec(m map[string]interface{}, spec FieldSpec) error {
	sep := "."
	paths := strings.Split(spec.Path, sep)
	curr := m
	for i, p := range paths {
		if i < len(paths)-1 {
			if _, ok := curr[p]; !ok {
				curr[p] = map[string]interface{}{}
			}
			candidate := curr[p]
			var ok bool
			curr, ok = candidate.(map[string]interface{})
			if !ok {
				return fmt.Errorf("%w: %#v path in conflict with value %#v",
					ErrCantSetField, spec.Path, candidate)
			}
		} else {
			curr[p] = spec.Value
		}
	}
	return nil
}

// AddField will add a field to the spec.
func (s *Spec) AddField(path string, val interface{}) {
	s.Fields = append(s.Fields, FieldSpec{
		Path: path, Value: val,
	})
}
