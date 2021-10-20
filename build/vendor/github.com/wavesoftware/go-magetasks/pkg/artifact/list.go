package artifact

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/wavesoftware/go-magetasks/config"
	"github.com/wavesoftware/go-magetasks/pkg/files"
	"github.com/wavesoftware/go-magetasks/pkg/output/color"
)

const (
	// ArtifactsBuilt is used to list artifacts that was built.
	ArtifactsBuilt     = "artifacts.built"
	allowReadAllOsPerm = 0o644
)

// ErrMisconfiguration when the project configuration faulty.
var ErrMisconfiguration = errors.New("project configuration is faulty")

// ListPublisher will output built artifacts as a "\n" delimited list in
// a result file.
type ListPublisher struct {
	FilePath         string
	ResultsRetriever func(config.Artifact) *config.Result
}

func (l ListPublisher) Accepts(artifact config.Artifact) bool {
	return l.resultsRetriever()(artifact) != nil
}

func (l ListPublisher) Publish(artifact config.Artifact, notifier config.Notifier) config.Result {
	reportFilename := l.FilePath
	if reportFilename == "" {
		reportFilename = "artifacts.list"
	}
	reportPath := path.Join(files.BuildDir(), reportFilename)
	result := l.resultsRetriever()(artifact)
	if result == nil && !result.Failed() {
		return config.Result{Error: fmt.Errorf(
			"%w: can't find result for %v", ErrMisconfiguration, artifact)}
	}
	artifactsList, ok := result.Info[ArtifactsBuilt].([]string)
	if !ok {
		return config.Result{Error: fmt.Errorf(
			"%w: can't find result for %v", ErrMisconfiguration, artifact)}
	}

	f, err := os.OpenFile(reportPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, allowReadAllOsPerm)
	defer func() {
		if f != nil {
			_ = f.Close()
		}
	}()
	if err != nil {
		return config.Result{Error: fmt.Errorf("%w: %v", ErrMisconfiguration, err)}
	}

	for _, artifactPath := range artifactsList {
		if _, err = f.WriteString(fmt.Sprintf("%s\n", artifactPath)); err != nil {
			return config.Result{Error: fmt.Errorf("%w: %v", ErrMisconfiguration, err)}
		}
	}

	notifier.Notify(fmt.Sprintf("artifact list: %s", color.Blue(reportPath)))

	return config.Result{}
}

func (l ListPublisher) resultsRetriever() func(config.Artifact) *config.Result {
	if l.ResultsRetriever != nil {
		return l.ResultsRetriever
	}
	return func(artifact config.Artifact) *config.Result {
		if !BinaryBuilder.Accepts(BinaryBuilder{}, artifact) {
			return nil
		}
		key := BuildKey(artifact)
		result, ok := config.Actual().Context.Value(key).(config.Result)
		if !ok {
			return nil
		}
		return &result
	}
}
