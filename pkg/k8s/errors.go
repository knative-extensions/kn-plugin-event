package k8s

import "errors"

var (
	// ErrInvalidReference if given reference is invalid.
	ErrInvalidReference = errors.New("reference is invalid")

	// ErrNotFound if given reference do not point to any resource.
	ErrNotFound = errors.New("resource not found")

	// ErrNotAddressable if found resource isn't addressable.
	ErrNotAddressable = errors.New("resource isn't addressable")

	// ErrMoreThenOneFound if more then one resource has been found.
	ErrMoreThenOneFound = errors.New("more then one resource has been found")

	// ErrUnexcpected if something unexpected actually has happened.
	ErrUnexcpected = errors.New("something unexpected actually has happened")

	// ErrICSenderJobFailed if the ICS job runner has failed.
	ErrICSenderJobFailed = errors.New("the ICS job runner has failed")
)
