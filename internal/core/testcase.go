package core

import (
	"errors"
	"fmt"
	"path/filepath"
	"slices"
	"strings"
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

func (tc *TestCase) Execute(fpath string) {
	ext := strings.TrimPrefix(filepath.Ext(fpath), ".")

	engine, ok := Engines[ext]
	if !ok {
		tc.Status = TestCaseStatusError
		tc.Details = fmt.Sprintf("no engine for extension: .%s", ext)
		return
	}

	output, err := engine.Run(fpath, tc.Input)
	tc.Output = output

	if err != nil {
		tc.Status = TestCaseStatusError
		switch {
		case errors.Is(err, ErrCompile):
			tc.Details = "Compilation Error"
		case errors.Is(err, ErrExecute):
			tc.Details = "Execution Error"
		default:
			tc.Details = err.Error()
		}
		return
	}

	tc.Details = "TODO: time taken"

	if strings.TrimRight(output, "\n") == strings.TrimRight(tc.Expected, "\n") {
		tc.Status = TestCaseStatusCorrect
	} else {
		tc.Status = TestCaseStatusWrong
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
