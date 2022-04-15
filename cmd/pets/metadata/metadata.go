package metadata

import "fmt"

// Build metadata. These values are intended to be overridden by values
// supplied by the build environment at link time. Values are specified
// by passing "-X package.Varname=Value" to the linker. For example:
//
//  go build -ldflags="-X meta.Name=BFF -X main.Version=1.0" ...etc...
//
// Ref: https://golang.org/cmd/link/
// nolint:gochecknoglobals // The flags must be global variables.
var (
	ApplicationName string
	CommitSha       string
	Version         string
)

// BuildMetadata is the metadata that can be populated for the build at runtime.
type BuildMetadata struct {
	ApplicationName string
	CommitSha       string
	Version         string
}

// NewBuildMetadata created a new BuildMetadata struct, populated with global values received at runtime.
func NewBuildMetadata() *BuildMetadata {
	return &BuildMetadata{
		ApplicationName: ApplicationName,
		CommitSha:       CommitSha,
		Version:         Version,
	}
}

// String converts the BuildMetadata struct to a string.
func (b *BuildMetadata) String() string {
	return fmt.Sprintf(
		"app=%s commitsha=%s version=%s",
		b.ApplicationName,
		b.CommitSha,
		b.Version,
	)
}
