package core

import "slices"

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

	Input    string `json:"input"`
	Expected string `json:"expected"`
	Output   string `json:"output"`
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
