package files

import (
	"os"
	"path"
	"strings"
	"sync"

	"github.com/magefile/mage/mg"
	"github.com/wavesoftware/go-ensure"
	"github.com/wavesoftware/go-magetasks/config"
	"github.com/wavesoftware/go-magetasks/pkg/dotenv"
	"github.com/wavesoftware/go-magetasks/pkg/output"
	"k8s.io/apimachinery/pkg/util/rand"
)

const randomDirLength = 12

var (
	randomDir    = "mage-build-" + rand.String(randomDirLength)
	buildDirOnce sync.Once
)

// EnsureBuildDir creates a build directory.
func EnsureBuildDir() {
	buildDirOnce.Do(func() {
		mg.Deps(dotenv.Load, output.Setup)
		d := path.Join(BuildDir(), "bin")
		ensure.NoError(os.MkdirAll(d, os.ModePerm))
		ensure.NoError(os.MkdirAll(ReportsDir(), os.ModePerm))
		if strings.Contains(ReportsDir(), randomDir) {
			output.Println("📁 Reports directory: ", ReportsDir())
		}
	})
}

// BuildDir returns project build dir.
func BuildDir() string {
	buildDir := os.Getenv("MAGE_BUILD_DIR")
	if buildDir != "" {
		return buildDir
	}
	buildDir = os.Getenv("BUILD_DIR")
	if buildDir != "" {
		return buildDir
	}
	return relativeTo(ProjectDir(), config.Actual().BuildDirPath...)
}

// ReportsDir returns project reports directory.
func ReportsDir() string {
	artifacts := os.Getenv("ARTIFACTS")
	if artifacts != "" {
		return path.Join(artifacts, randomDir)
	}
	return relativeTo(BuildDir(), "reports")
}

// ProjectDir returns project repo directory.
func ProjectDir() string {
	if config.Actual().ProjectDir != "" {
		return config.Actual().ProjectDir
	}
	repoDir, err := os.Getwd()
	ensure.NoError(err)
	return repoDir
}

func relativeTo(to string, paths ...string) string {
	fullpath := make([]string, len(paths)+1)
	fullpath[0] = to
	for ix, elem := range paths {
		fullpath[ix+1] = elem
	}
	return path.Join(fullpath...)
}
