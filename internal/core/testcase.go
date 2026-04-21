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

func LoadTestCases(filePath string) (TestCaseList, error) {
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

		sample := TestCaseList{
			&TestCase{
				Status:  TestCaseStatusPending,
				Details: "Idle",
			},
		}
		if err := sample.Save(testCaseFile); err != nil {
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

func (list *TestCaseList) Create() {
	*list = append(*list, &TestCase{
		Status:  TestCaseStatusPending,
		Details: "Idle",
	})
}

func (list *TestCaseList) Delete(index int) {
	*list = slices.Delete(*list, index, index+1)
}
