package logging

import (
	"encoding/json"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// JSON returns a zap.Field that dumps an object as JSON.
func JSON(name string, o interface{}) zap.Field {
	return zap.Stringer(name, jsonMarshal(func() string {
		bytes, err := json.Marshal(o)
		if err != nil {
			panic(errors.WithStack(err))
		}
		return string(bytes)
	}))
}

type jsonMarshal func() string

func (jm jsonMarshal) String() string {
	return jm()
}
