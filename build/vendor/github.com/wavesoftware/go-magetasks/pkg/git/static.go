package git

import "github.com/wavesoftware/go-magetasks/pkg/strings"

// StaticRepository is stub repo implementation which gives back values provided
// to it upfront.
type StaticRepository struct {
	DescribeString string
	TagsSet        strings.Set
}

func (s StaticRepository) Describe() (string, error) {
	return s.DescribeString, nil
}

func (s StaticRepository) Tags() ([]string, error) {
	return s.TagsSet.Slice(), nil
}
