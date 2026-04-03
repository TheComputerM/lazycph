package core

type TestCaseStatus string

const (
	TestCaseStatusCorrect TestCaseStatus = "Correct"
	TestCaseStatusPending TestCaseStatus = "Pending"
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

func GetTestCases() ([]TestCase, error) {
	return []TestCase{
		{TestCaseStatusCorrect, "200ms", "STDIN:1", "EXPECTED:1", "STDOUT:1"},
		{TestCaseStatusPending, "Queued", "STDIN:pending", "EXPECTED:pending", "STDOUT:"},
		{TestCaseStatusError, "Compilation Error", "STDIN:2", "EXPECTED:2", "STDOUT:2"},
		{TestCaseStatusWrong, "300ms", "STDIN:3", "EXPECTED:3", "STDOUT:3"},
	}, nil
}

func CreateTestCase(list *[]TestCase) {
	*list = append(*list, TestCase{
		Status:   TestCaseStatusPending,
		Details:  "Idle",
		Input:    "",
		Expected: "",
		Output:   "",
	})
}

func RunTestCase(testCase *TestCase) {
	testCase.Status = TestCaseStatusCorrect
	testCase.Details = "300ms"
	testCase.Output = "Some text"
}
