package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/cardil/kn-event/internal/event"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/ghodss/yaml"
)

var (
	// ErrUnsupportedOutputMode if user passed unsupported output mode.
	ErrUnsupportedOutputMode = errors.New("unsupported output mode")

	// ErrInvalidFormat if user pass an un-parsable format.
	ErrInvalidFormat = errors.New("invalid format")
)

// CreateWithArgs will create an event by parsing given args.
func (c *App) CreateWithArgs(args *EventArgs) (*cloudevents.Event, error) {
	spec := &event.Spec{
		Type:   args.Type,
		ID:     args.ID,
		Source: args.Source,
		Fields: make([]event.FieldSpec, 0, len(args.Fields)+len(args.RawFields)),
	}
	for _, fieldAssigment := range args.Fields {
		splitted := strings.SplitN(fieldAssigment, "=", 2)
		path, value := splitted[0], splitted[1]
		if boolVal, err := readAsBoolean(value); err == nil {
			spec.AddField(path, boolVal)
		} else if floatVal, err := readAsFloat64(value); err == nil {
			spec.AddField(path, floatVal)
		} else {
			spec.AddField(path, value)
		}
	}
	for _, fieldAssigment := range args.RawFields {
		splitted := strings.SplitN(fieldAssigment, "=", 2)
		path, value := splitted[0], splitted[1]
		spec.AddField(path, value)
	}
	return event.CreateFromSpec(spec)
}

// PresentWith will present an event with specified output.
func (c *App) PresentWith(e *cloudevents.Event, output OutputMode) (string, error) {
	switch output {
	case HumanReadable:
		return presentEventAsHumanReadable(e)
	case JSON:
		return presentEventAsJSON(e)
	case YAML:
		return presentEventAsYaml(e)
	}
	return "", fmt.Errorf("%w: %v", ErrUnsupportedOutputMode, output)
}

func presentEventAsYaml(in *cloudevents.Event) (string, error) {
	j, err := presentEventAsJSON(in)
	if err != nil {
		return "", err
	}
	e := cloudevents.Event{}
	err = json.Unmarshal([]byte(j), &e)
	if err != nil {
		return "", err
	}
	bytes, err := yaml.Marshal(e)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func presentEventAsJSON(event *cloudevents.Event) (string, error) {
	bytes, err := json.MarshalIndent(event, "", "  ")
	var out string
	if err == nil {
		out = string(bytes)
	}
	return out, err
}

func presentEventAsHumanReadable(e *cloudevents.Event) (string, error) {
	formattedTime := e.Time().
		In(time.UTC).
		Format(time.RFC3339Nano)
	m := map[string]interface{}{}
	err := json.Unmarshal(e.Data(), &m)
	if err != nil {
		return "", err
	}
	data, err := json.MarshalIndent(m, "  ", "  ")
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(
		`☁️  cloudevents.Event
Validation: valid
Context Attributes,
  specversion: %s
  type: %s
  source: %s
  id: %s
  time: %s
  datacontenttype: %s
Data,
  %s`,
		e.SpecVersion(),
		e.Type(),
		e.Source(),
		e.ID(),
		formattedTime,
		e.DataContentType(),
		string(data),
	), nil
}

func readAsBoolean(in string) (bool, error) {
	val, err := strconv.ParseBool(in)
	// TODO(cardil): log error as it may be beneficial for debugging
	if err != nil {
		return false, err
	}
	text := fmt.Sprintf("%t", val)
	if in == text {
		return val, nil
	}
	return false, fmt.Errorf("%w: not a bool: %v", ErrInvalidFormat, in)
}

func readAsFloat64(in string) (float64, error) {
	if intVal, err := readAsInt64(in); err == nil {
		return float64(intVal), err
	}
	val, err := strconv.ParseFloat(in, 64)
	// TODO(cardil): log error as it may be beneficial for debugging
	if err != nil {
		return -0, err
	}
	text := fmt.Sprintf("%f", val)
	if in == text {
		return val, nil
	}
	return -0, fmt.Errorf("%w: not a float: %v", ErrInvalidFormat, in)
}

func readAsInt64(in string) (int64, error) {
	val, err := strconv.ParseInt(in, 10, 64)
	// TODO(cardil): log error as it may be beneficial for debugging
	if err != nil {
		return -0, err
	}
	text := fmt.Sprintf("%d", val)
	if in == text {
		return val, nil
	}
	return -0, fmt.Errorf("%w: not an int: %v", ErrInvalidFormat, in)
}
