package extract

import (
	"fmt"

	"github.com/pkg/errors"
)

var (
	// ErrBug is an error that indicates a bug in the code.
	ErrBug = errors.New("probably a bug in the code")

	// ErrUnexpected is an error that indicates an unexpected situation.
	ErrUnexpected = errors.New("unexpected situation")
)

func wrapErr(err error, target error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, target) {
		return err
	}
	return errors.WithStack(fmt.Errorf("%w: %v", target, err))
}
