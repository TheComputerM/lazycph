package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "charm.land/bubbletea/v2"
)

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
