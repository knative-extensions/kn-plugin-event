package tests

import (
	"encoding/json"
	"errors"
	"fmt"
)

// ErrCantUnmarshalData if given bytes can't be unmarshalled as JSON.
var ErrCantUnmarshalData = errors.New("can't unmarshal data")

// UnmarshalCloudEventData will take bytes and unmarshall it as JSON to map structure.
func UnmarshalCloudEventData(bytes []byte) (map[string]interface{}, error) {
	m := map[string]interface{}{}
	err := json.Unmarshal(bytes, &m)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCantUnmarshalData, err)
	}
	return m, nil
}
