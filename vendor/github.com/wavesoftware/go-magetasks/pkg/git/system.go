package git

import (
	"regexp"
	"strings"

	"github.com/magefile/mage/sh"
)

// ref.: https://regex101.com/r/ppnq02/1
var remoteTagPattern = regexp.MustCompile(`^[0-9a-f]{7,}\s+refs/tags/([^^]+)(:?\^{})$`)

type installedGitBinaryRepo struct {
	Remote
}

func (s installedGitBinaryRepo) Describe() (string, error) {
	return sh.Output("git", "describe", "--always", "--tags", "--dirty")
}

func (s installedGitBinaryRepo) Tags() ([]string, error) {
	output, err := sh.Output("git", "ls-remote", "--tags", s.remote())
	if err != nil {
		return nil, err
	}
	return parseLsRemoteTagsOutput(output), nil
}

func parseLsRemoteTagsOutput(output string) []string {
	lines := strings.Split(output, "\n")
	tagsMap := make(map[string]bool)
	for _, line := range lines {
		match := remoteTagPattern.FindSubmatch([]byte(line))
		if match == nil {
			continue
		}
		tag := string(match[1])
		tagsMap[tag] = true
	}
	tags := make([]string, 0, len(tagsMap))
	for tag := range tagsMap {
		tags = append(tags, tag)
	}
	return tags
}

func (s installedGitBinaryRepo) remote() string {
	if s.Remote.URL != "" {
		return s.Remote.URL
	}
	return s.Remote.Name
}
