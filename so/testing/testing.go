package testing

import (
	"solod.dev/so/fmt"
	"solod.dev/so/io"
	"solod.dev/so/os"
)

// T is the context passed to a test function. It records failure and skip
// state for a single test.
//
// Unlike Go's testing.T, the message methods are NOT variadic: So cannot
// forward a "...any" argument to fmt, so there is no Errorf/Fatalf. Format the
// message yourself with fmt.Sprintf when needed:
//
//	buf := fmt.NewBuffer(64)
//	t.Error(fmt.Sprintf(buf, "Index = %d, want 6", got))
//
// So also has no recover, so T cannot unwind a running test. Fatal only marks
// the test failed and prints the message; by convention the test function must
// return right after calling it:
//
//	if err != nil {
//		t.Fatal("open failed")
//		return
//	}
//
// By design, a hard crash (panic or segfault) in a test aborts
// the entire test run, not just the current test.
type T struct {
	name    string
	w       io.Writer
	failed  bool
	skipped bool
}

// Name returns the name of the running test.
func (t *T) Name() string { return t.name }

// Failed reports whether the test has failed.
func (t *T) Failed() bool { return t.failed }

// Fail marks the test as failed but continues execution.
func (t *T) Fail() { t.failed = true }

// Log records msg in the test log.
func (t *T) Log(msg string) {
	fmt.Fprintf(t.w, "    %s\n", msg)
}

// Error is equivalent to Log followed by Fail.
func (t *T) Error(msg string) {
	t.Log(msg)
	t.Fail()
}

// Fatal is equivalent to Log followed by Fail. The test function must return
// after calling it; see [T].
func (t *T) Fatal(msg string) {
	t.Log(msg)
	t.Fail()
}

// Skip marks the test as skipped. Like Fatal, the test must return afterwards.
func (t *T) Skip(msg string) {
	t.Log(msg)
	t.skipped = true
}

// Test represents a single test to be run by the test runner.
type Test struct {
	Name string
	F    func(t *T)
}

// RunTests runs the given tests for package pkg, prints per-test results
// to stdout, and exits with a non-zero status if any test failed.
func RunTests(pkg string, tests []Test) {
	failed := 0
	skipped := 0
	for _, tc := range tests {
		t := &T{name: tc.Name, w: os.Stdout}
		fmt.Fprintf(t.w, "=== RUN   %s\n", t.name)
		tc.F(t)

		if t.skipped {
			fmt.Fprintf(t.w, "--- SKIP: %s\n", t.name)
			skipped++
			continue
		}
		if t.failed {
			fmt.Fprintf(t.w, "--- FAIL: %s\n", t.name)
			failed++
			continue
		}
		fmt.Fprintf(t.w, "--- PASS: %s\n", t.name)
	}

	if failed > 0 {
		fmt.Fprintf(os.Stdout, "FAIL\t%s\t%d of %d failed\n", pkg, failed, len(tests))
		os.Exit(1)
	}
	if skipped > 0 {
		fmt.Fprintf(os.Stdout, "ok\t%s\t%d tests (%d skipped)\n", pkg, len(tests), skipped)
		return
	}
	fmt.Fprintf(os.Stdout, "ok\t%s\t%d tests\n", pkg, len(tests))
}
