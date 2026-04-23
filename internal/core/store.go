package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// NewTestCase returns a TestCase in its initial idle state.
func NewTestCase() *TestCase {
	return &TestCase{
		Status:  TestCaseStatusPending,
		Details: "Idle",
	}
}

func LoadTestCaseList(filePath string) TestCaseList {
	list, err := loadTestCaseList(filePath)

	if err != nil || len(list) == 0 {
		// return a new test case list if error with existing one
		return TestCaseList{NewTestCase()}
	}

	return list
}

func loadTestCaseList(filePath string) (TestCaseList, error) {
	lazyCphDir := filepath.Join(filepath.Dir(filePath), ".lazycph")
	storeFile := filepath.Join(lazyCphDir, strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))+".json")

	data, err := os.ReadFile(storeFile)
	if err != nil {
		return nil, err
	}

	var list TestCaseList
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return nil, fmt.Errorf("testcase list is empty")
	}

	return list, nil
}

func (list TestCaseList) Save(sourceFile string) error {
	lazyCphDir := filepath.Join(filepath.Dir(sourceFile), ".lazycph")
	if info, err := os.Stat(lazyCphDir); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}

		if err := os.Mkdir(lazyCphDir, 0o755); err != nil {
			return err
		}

		if err := os.WriteFile(filepath.Join(lazyCphDir, ".gitignore"), []byte("*"), 0o644); err != nil {
			return err
		}
	} else if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", lazyCphDir)
	}

	data, err := json.Marshal(list)
	if err != nil {
		return err
	}

	storeFile := filepath.Join(lazyCphDir, strings.TrimSuffix(filepath.Base(sourceFile), filepath.Ext(sourceFile))+".json")

	return os.WriteFile(storeFile, data, 0o644)
}


