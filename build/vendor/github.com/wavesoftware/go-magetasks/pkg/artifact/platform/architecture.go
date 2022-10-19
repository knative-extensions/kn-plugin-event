package platform

// Architecture is an CPU architecture.
type Architecture string

const (
	AMD64   Architecture = "amd64"
	ARM64   Architecture = "arm64"
	S390X   Architecture = "s390x"
	PPC64LE Architecture = "ppc64le"
)
