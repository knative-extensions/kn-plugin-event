package git

// Repository gives information about the git SCM repository.
type Repository interface {
	// Describe will return an output that closely matches that of
	// "git describe --tags --always --dirty" for given repository.
	Describe() (string, error)
	// Tags will return repository tags.
	Tags() ([]string, error)
}
