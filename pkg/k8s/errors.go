package k8s

import "errors"

var (
	// ErrInvalidReference if given reference is invalid.
	ErrInvalidReference = errors.New("reference is invalid")

	// ErrNotAddressable if found resource isn't addressable.
	ErrNotAddressable = errors.New("resource isn't addressable")

	// ErrUnexcpected if something unexpected actually has happened.
	ErrUnexcpected = errors.New("something unexpected actually has happened")

	// ErrJobFailed if the job has failed.
	ErrJobFailed = errors.New("the Kubernetes job failed")
)
