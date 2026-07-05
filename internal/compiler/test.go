package compiler

import (
	"fmt"
	"strings"
)

// testKind discovers TestXxx(t *testing.T) functions
// and runs them via testing.RunTests.
var testKind = kind{
	subdir:  "test",
	command: "so test",
	noun:    "test",
	prefix:  "Test",
	param:   "T",
	emit:    emitTests,
}

// Test discovers TestXxx functions in the "test" subdirectory of srcDir,
// generates a deterministic main.go runner there, and runs it via Run.
func Test(srcDir string, opts Options) error {
	return testKind.run(srcDir, opts)
}

// emitTests writes the runner body dispatching the tests via testing.RunTests.
func emitTests(b *strings.Builder, pkg string, names []string) {
	b.WriteString("import \"solod.dev/so/testing\"\n\n")
	b.WriteString("func main() {\n")
	fmt.Fprintf(b, "\ttesting.RunTests(%q, []testing.Test{\n", pkg)
	for _, name := range names {
		fmt.Fprintf(b, "\t\t{Name: %q, F: %s},\n", name, name)
	}
	b.WriteString("\t})\n")
	b.WriteString("}\n")
}
