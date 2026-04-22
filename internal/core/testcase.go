package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
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

// TestCaseExecutedMsg is dispatched when a TestCase has finished running.
type TestCaseExecutedMsg struct {
	TestCase *TestCase
}

// Execute runs the test case against the source file at srcPath.
func (tc *TestCase) ExecuteCmd(srcPath string) tea.Cmd {
	tc.Status = TestCaseStatusPending
	tc.Details = "Running..."

	return func() tea.Msg {
		tc.Execute(srcPath)
		return TestCaseExecutedMsg{TestCase: tc}
	}
}

// run performs the synchronous execution and updates tc with the verdict.
func (tc *TestCase) Execute(srcPath string) {
	ext := strings.TrimPrefix(filepath.Ext(srcPath), ".")
	engine, ok := Engines[ext]
	if !ok {
		tc.Status = TestCaseStatusError
		tc.Details = "No Engine"
		tc.Output = fmt.Sprintf("LazyCPH has no execution engine for .%s, add one at ~/.config/lazycph.json", ext)
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

// newTestCase returns a TestCase in its initial idle state.
func newTestCase() *TestCase {
	return &TestCase{
		Status:  TestCaseStatusPending,
		Details: "Idle",
	}
}

// Append adds a fresh, idle test case to the end of the list.
func (list *TestCaseList) Append() {
	*list = append(*list, newTestCase())
}

// RemoveAt removes the test case at index. It is a no-op if index is out of
// range.
func (list *TestCaseList) RemoveAt(index int) {
	if index < 0 || index >= len(*list) {
		return
	}
	*list = slices.Delete(*list, index, index+1)
}

func LoadTestCaseList(filePath string) (TestCaseList, error) {
	if err := ensureTestCaseFile(filePath); err != nil {
		return nil, err
	}

	lazyCphDir := filepath.Join(filepath.Dir(filePath), ".lazycph")
	testCaseFile := filepath.Join(lazyCphDir, strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))+".json")

	data, err := os.ReadFile(testCaseFile)
	if err != nil {
		return nil, err
	}

	var list TestCaseList
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, err
	}

	return list, nil
}

func ensureTestCaseFile(filePath string) error {
	lazyCphDir := filepath.Join(filepath.Dir(filePath), ".lazycph")
	if _, err := os.Stat(lazyCphDir); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}

		if err := os.Mkdir(lazyCphDir, 0o755); err != nil {
			return err
		}

		if err := os.WriteFile(filepath.Join(lazyCphDir, ".gitignore"), []byte("*"), 0o644); err != nil {
			return err
		}
	}

	testCaseFile := filepath.Join(lazyCphDir, strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))+".json")
	if _, err := os.Stat(testCaseFile); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}

		sample := TestCaseList{newTestCase()}
		if err := sample.Save(filePath); err != nil {
			return err
		}
	}

	return nil
}

func (list TestCaseList) Save(filePath string) error {
	// assume lazycph dir exists
	lazyCphDir := filepath.Join(filepath.Dir(filePath), ".lazycph")
	testCaseFile := filepath.Join(lazyCphDir, strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))+".json")

	data, err := json.Marshal(list)
	if err != nil {
		return err
	}

	return os.WriteFile(testCaseFile, data, 0o644)
}

func (list TestCaseList) SaveCmd(filePath string) tea.Cmd {
	return func() tea.Msg {
		return list.Save(filePath)
	}
}
