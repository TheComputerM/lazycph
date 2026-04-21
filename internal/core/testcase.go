package core

import (
	"errors"
	"fmt"
	"path/filepath"
	"slices"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
)

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

	// STDIN for the test case
	Input string `json:"input"`

	// Expected STDOUT for the test case
	Expected string `json:"expected"`

	// Actual STDOUT from the test case
	Output string `json:"output"`
}

type TestCaseExecutedMsg struct {
	TestCase *TestCase
}

func (tc *TestCase) Execute(fpath string) tea.Cmd {
	tc.Status = TestCaseStatusPending
	tc.Details = "Running..."

	return func() tea.Msg {
		ext := strings.TrimPrefix(filepath.Ext(fpath), ".")

		engine, ok := Engines[ext]
		if !ok {
			tc.Status = TestCaseStatusError
			tc.Details = fmt.Sprintf("no engine for extension: .%s", ext)
			return TestCaseExecutedMsg{TestCase: tc}
		}

		start := time.Now()
		output, err := engine.Run(fpath, tc.Input)
		elapsed := time.Since(start)
		tc.Output = output

		if err != nil {
			tc.Status = TestCaseStatusError
			switch {
			case errors.Is(err, ErrCompile):
				tc.Details = "Compilation Error"
			case errors.Is(err, ErrExecute):
				tc.Details = "Execution Error"
			}
			return TestCaseExecutedMsg{TestCase: tc}
		}

		tc.Details = fmt.Sprintf("%dms", elapsed.Milliseconds())

		if strings.TrimRight(output, "\n") == strings.TrimRight(tc.Expected, "\n") {
			tc.Status = TestCaseStatusCorrect
		} else {
			tc.Status = TestCaseStatusWrong
		}
		return TestCaseExecutedMsg{TestCase: tc}
	}
}

type TestCaseList []*TestCase

func GetTestCases() (TestCaseList, error) {
	// mock data
	return TestCaseList{
		{TestCaseStatusCorrect, "200ms", "STDIN:1", "EXPECTED:1", "STDOUT:1"},
		{TestCaseStatusPending, "Queued", "STDIN:pending", "EXPECTED:pending", "STDOUT:"},
		{TestCaseStatusError, "Compilation Error", "STDIN:2", "EXPECTED:2", "STDOUT:2"},
		{TestCaseStatusWrong, "300ms", "STDIN:3", "EXPECTED:3", "STDOUT:3"},
	}, nil
}

func (list *TestCaseList) Create() {
	*list = append(*list, &TestCase{
		Status:  TestCaseStatusPending,
		Details: "Idle",
	})
}

func (list *TestCaseList) Delete(index int) {
	*list = slices.Delete(*list, index, index+1)
}
