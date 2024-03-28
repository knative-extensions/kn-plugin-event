package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/ghodss/yaml"
	"knative.dev/kn-plugin-event/pkg/event"
)

const (
	fieldAssigmentSize = 2
	decimalBase        = 10
	precision64BitSize = 64
)

var (
	// ErrUnsupportedOutputMode if user passed unsupported output mode.
	ErrUnsupportedOutputMode = errors.New("unsupported output mode")

	// ErrInvalidFormat if user pass an un-parsable format.
	ErrInvalidFormat = errors.New("invalid format")

	// ErrCantBuildEvent if event can't be built.
	ErrCantBuildEvent = errors.New("can't build event")

	// ErrCantMarshalEvent if event can't be marshalled to text.
	ErrCantMarshalEvent = errors.New("can't marshal event")
)

// CreateWithArgs will create an event by parsing given args.
func (a *App) CreateWithArgs(args *EventArgs) (*cloudevents.Event, error) {
	spec := &event.Spec{
		Type:   args.Type,
		ID:     args.ID,
		Source: args.Source,
		Fields: make([]event.FieldSpec, 0, len(args.Fields)+len(args.RawFields)),
	}
	for _, fieldAssigment := range args.Fields {
		split := strings.SplitN(fieldAssigment, "=", fieldAssigmentSize)
		path, value := split[0], split[1]
		var floatVal float64
		if boolVal, err := readAsBoolean(value); err == nil {
			spec.AddField(path, boolVal)
		} else if floatVal, err = readAsFloat64(value); err == nil {
			spec.AddField(path, floatVal)
		} else {
			spec.AddField(path, value)
		}
	}
	for _, fieldAssigment := range args.RawFields {
		split := strings.SplitN(fieldAssigment, "=", fieldAssigmentSize)
		path, value := split[0], split[1]
		spec.AddField(path, value)
	}
	ce, err := event.CreateFromSpec(spec)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrCantBuildEvent, err)
	}
	return ce, nil
}

// PresentWith will present an event with specified output.
func (a *App) PresentWith(e *cloudevents.Event, output OutputMode) (string, error) {
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
	bytes, err := yaml.Marshal(in)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrCantMarshalEvent, err)
	}
	return string(bytes), nil
}

func presentEventAsJSON(event *cloudevents.Event) (string, error) {
	bytes, err := json.MarshalIndent(event, "", "  ")
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrCantMarshalEvent, err)
	}
	return string(bytes), nil
}

func presentEventAsHumanReadable(e *cloudevents.Event) (string, error) {
	formattedTime := e.Time().
		In(time.UTC).
		Format(time.RFC3339Nano)
	m := map[string]interface{}{}
	err := json.Unmarshal(e.Data(), &m)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrCantMarshalEvent, err)
	}
	data, err := json.MarshalIndent(m, "  ", "  ")
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrCantMarshalEvent, err)
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
		return false, fmt.Errorf("%w: %w", ErrInvalidFormat, err)
	}
	if text := strconv.FormatBool(val); in == text {
		return val, nil
	}
	return false, fmt.Errorf("%w: not a bool: %v", ErrInvalidFormat, in)
}

func readAsFloat64(in string) (float64, error) {
	if intVal, err := readAsInt64(in); err == nil {
		return float64(intVal), nil
	}
	val, err := strconv.ParseFloat(in, precision64BitSize)
	// TODO(cardil): log error as it may be beneficial for debugging
	if err != nil {
		return -0, fmt.Errorf("%w: %w", ErrInvalidFormat, err)
	}
	if text := fmt.Sprintf("%f", val); in == text {
		return val, nil
	}
	return -0, fmt.Errorf("%w: not a float: %v", ErrInvalidFormat, in)
}

func readAsInt64(in string) (int64, error) {
	val, err := strconv.ParseInt(in, decimalBase, precision64BitSize)
	// TODO(cardil): log error as it may be beneficial for debugging
	if err != nil {
		return -0, fmt.Errorf("%w: %w", ErrInvalidFormat, err)
	}
	if text := strconv.FormatInt(val, 10); in == text {
		return val, nil
	}
	return -0, fmt.Errorf("%w: not an int: %v", ErrInvalidFormat, in)
}
