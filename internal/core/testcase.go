package core

import (
	"errors"
	"fmt"
	"path/filepath"
	"slices"
	"strings"
	"time"
)

// TestCaseStatus is the verdict of the most recent run of a TestCase.
type TestCaseStatus string

const (
	TestCaseStatusPending TestCaseStatus = "Pending"
	TestCaseStatusCorrect TestCaseStatus = "Correct"
	TestCaseStatusWrong   TestCaseStatus = "Wrong"
	TestCaseStatusError   TestCaseStatus = "Error"
)

type TestCase struct {
	Status  TestCaseStatus `json:"status"`
	Details string         `json:"details"`

	// Input is the data fed to the program's stdin.
	Input string `json:"input"`

	// Expected is the stdout the program should produce.
	Expected string `json:"expected"`

	// Output is the actual stdout (and stderr) captured during the last run.
	Output string `json:"output"`
}

// run performs the synchronous execution and updates tc with the verdict.
func (tc *TestCase) Execute(srcPath string) {
	ext := filepath.Ext(srcPath)
	engine, ok := Engines[ext]
	if !ok {
		tc.Status = TestCaseStatusError
		tc.Details = "No Engine"
		tc.Output = fmt.Sprintf("LazyCPH has no execution engine for %s, add one at ~/.config/lazycph.json", ext)
		return
	}

	start := time.Now()
	output, err := engine.Run(srcPath, tc.Input)
	tc.Output = output

	if err != nil {
		tc.Status = TestCaseStatusError

		switch {
		case errors.Is(err, ErrCompile):
			tc.Details = "Compilation Error"
		case errors.Is(err, ErrExecute):
			tc.Details = "Runtime Error"
		default:
			tc.Details = "Unknown Error"
			tc.Output = err.Error()
		}
		return
	}

	tc.Details = fmt.Sprintf("%dms", time.Since(start).Milliseconds())
	if outputMatches(output, tc.Expected) {
		tc.Status = TestCaseStatusCorrect
	} else {
		tc.Status = TestCaseStatusWrong
	}
}

// outputMatches reports whether actual matches expected, ignoring trailing
// newlines on either side.
func outputMatches(actual, expected string) bool {
	return strings.TrimRight(actual, "\n") == strings.TrimRight(expected, "\n")
}

// TestCaseList is an ordered collection of test cases for a single source file.
type TestCaseList []*TestCase

// Append adds a fresh, idle test case to the end of the list.
func (list *TestCaseList) Append() {
	*list = append(*list, NewTestCase())
}

// RemoveAt removes the test case at index. It is a no-op if index is out of
// range.
func (list *TestCaseList) RemoveAt(index int) {
	if index < 0 || index >= len(*list) {
		return
	}
	*list = slices.Delete(*list, index, index+1)
}
