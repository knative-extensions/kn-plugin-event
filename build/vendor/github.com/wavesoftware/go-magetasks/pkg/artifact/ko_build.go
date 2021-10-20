package artifact

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/ko/pkg/build"
	"github.com/google/ko/pkg/commands"
	"github.com/google/ko/pkg/commands/options"
	"github.com/wavesoftware/go-magetasks/config"
	"github.com/wavesoftware/go-magetasks/pkg/files"
	"github.com/wavesoftware/go-magetasks/pkg/ldflags"
	"github.com/wavesoftware/go-magetasks/pkg/output/color"
	"golang.org/x/mod/modfile"
)

const (
	koImportPath  = "ko.import.path"
	koBuildResult = "ko.build.result"
)

// ErrKoFailed when th Google's ko fails to build.
var ErrKoFailed = errors.New("ko failed")

// KoBuilder builds images with Google's KO.
type KoBuilder struct{}

func (kb KoBuilder) Accepts(artifact config.Artifact) bool {
	_, ok := artifact.(Image)
	return ok
}

func (kb KoBuilder) Build(artifact config.Artifact, notifier config.Notifier) config.Result {
	image, ok := artifact.(Image)
	if !ok {
		return config.Result{Error: ErrInvalidArtifact}
	}
	importPath, err := imageImportPath(image)
	if err != nil {
		return resultErrKoFailed(err)
	}
	bo := &options.BuildOptions{
		Platform: buildPlatformString(image),
		Labels:   buildLabels(image, importPath),
	}
	fillInLdflags(bo, importPath, image)
	ctx := config.Actual().Context
	builder, err := commands.NewBuilder(ctx, bo)
	if err != nil {
		return resultErrKoFailed(err)
	}
	result, err := builder.Build(ctx, importPath)
	if err != nil {
		return resultErrKoFailed(err)
	}
	ref, err := calculateImageReference(result, artifact)
	if err != nil {
		return resultErrKoFailed(err)
	}
	notifier.Notify(fmt.Sprintf("ko built image: %s", color.Blue(ref)))
	return config.Result{Info: map[string]interface{}{
		imageReferenceKey: ref.String(),
		koBuildResult:     result,
	}}
}

func fillInLdflags(bo *options.BuildOptions, importPath string, image Image) {
	c := config.Actual()
	args := make([]string, 0)
	if c.Version != nil || len(c.BuildVariables) > 0 || len(image.BuildVariables) > 0 {
		builder := ldflags.NewBuilder()
		for key, resolver := range c.BuildVariables {
			builder.Add(key, resolver)
		}
		if c.Version != nil {
			builder.Add(c.Version.Path, c.Version.Resolver)
		}
		for key, resolver := range image.BuildVariables {
			builder.Add(key, resolver)
		}
		args = builder.Build()
	}
	if len(args) > 0 {
		bo.BuildConfigs = map[string]build.Config{
			importPath: {
				ID:      "ldflags-config",
				Ldflags: args,
			},
		}
	}
}

func buildLabels(image Image, importPath string) []string {
	labels := make([]string, 0, len(image.Labels))
	if version := config.Actual().Version; version != nil {
		labels = append(labels, fmt.Sprintf("version=%s", version.Resolver()))
	}
	labels = append(labels, fmt.Sprintf("%s=%s", koImportPath, importPath))
	for key, resolver := range image.Labels {
		labels = append(labels, fmt.Sprintf("%s=%s", key, resolver()))
	}
	return labels
}

func buildPlatformString(im Image) string {
	platforms := make([]string, len(im.Architectures))
	for i, architecture := range im.Architectures {
		platforms[i] = fmt.Sprintf("linux/%s", architecture)
	}
	return strings.Join(platforms, ",")
}

func calculateImageReference(result build.Result, artifact config.Artifact) (name.Reference, error) {
	kp := KoPublisher{}
	po, err := kp.publishOptions()
	if err != nil {
		return nil, err
	}
	po.Push = false
	po.TarballFile = ""
	po.OCILayoutPath = ""
	po.Local = false
	publisher, err := commands.NewPublisher(po)
	if err != nil {
		return nil, err
	}
	defer closePublisher(publisher)
	ctx := config.Actual().Context
	ref, err := publisher.Publish(ctx, result, artifact.GetName())
	if err != nil {
		return nil, err
	}
	return ref, nil
}

func resultErrKoFailed(err error) config.Result {
	return config.Result{
		Error: fmt.Errorf("%w: %v", ErrKoFailed, err),
	}
}

func imageImportPath(image Image) (string, error) {
	binDir := fullBinaryDirectory(image.GetName())
	rs, err := lookForGoModule(binDir)
	if err != nil {
		return "", err
	}
	importPath := rs.resolve(binDir)
	if resolver, ok := image.Labels[koImportPath]; ok {
		importPath = resolver()
	}
	return importPath, nil
}

func lookForGoModule(dir string) (lookupGoModuleResult, error) {
	rs := lookupGoModuleResult{}
	for i := 0; i < 10_000; i++ {
		modFile := path.Join(dir, "go.mod")
		if files.DontExists(modFile) {
			dir = path.Dir(dir)
			rs.directoryDistance++
			continue
		}
		bytes, err := ioutil.ReadFile(modFile)
		if err != nil {
			return rs, err
		}
		file, err := modfile.Parse(modFile, bytes, nil)
		if err != nil {
			return rs, err
		}
		rs.module = file
		return rs, nil
	}
	return rs, fmt.Errorf("%w: can't find go module", ErrKoFailed)
}

type lookupGoModuleResult struct {
	module            *modfile.File
	directoryDistance int
}

func (r lookupGoModuleResult) resolve(dir string) string {
	root := dir
	for i := 0; i < r.directoryDistance; i++ {
		root = path.Dir(root)
	}
	p := strings.Replace(dir, root, "", 1)
	return path.Join(r.module.Module.Mod.Path, p)
}
