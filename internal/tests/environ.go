package tests

import (
	"fmt"
	"os"
)

// WithEnviron will execute a block of code with temporal environment set.
func WithEnviron(env map[string]string, body func() error) error {
	old := map[string]*string{}
	for k := range env {
		var vv *string
		if v, ok := os.LookupEnv(k); ok {
			vv = &v
		}
		old[k] = vv
	}
	for k, v := range env {
		err := os.Setenv(k, v)
		if err != nil {
			return fmt.Errorf("can't set env %v=%v: %w", k, v, err)
		}
	}
	defer func() {
		for k, v := range old {
			if v != nil {
				_ = os.Setenv(k, *v)
			} else {
				_ = os.Unsetenv(k)
			}
		}
	}()
	return body()
}
