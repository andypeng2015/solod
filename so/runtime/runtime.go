// Package runtime provides information about
// the environment where the program was compiled.
package runtime

//so:embed runtime.h
var runtime_h string

//so:extern runtime_buildVersion
var buildVersion string

// Version returns the So tree's version string.
// It is either the commit hash of the build or,
// when possible, a release tag like "v0.1.0".
func Version() string {
	return buildVersion
}

// GOOS is the running program's operating system target:
// one of darwin, linux, windows, and so on.
//
//so:extern
const GOOS string = "unknown"

// GOARCH is the running program's architecture target:
// one of amd64, arm64, and so on.
//
//so:extern
const GOARCH string = "unknown"

// Recognized GOOS/GOARCH pairs are:
// GOOS		GOARCH
// darwin	amd64
// darwin	arm64
// linux	amd64
// linux	arm64
// windows	amd64
// windows	arm64
