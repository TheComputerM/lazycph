package core

import (
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

func testDataPath(rel string) string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "../..", "test", "data", rel)
}

func TestEngine(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		expected := "Hello World\n"
		t.Run("cpp", func(t *testing.T) {
			if _, err := exec.LookPath("c++"); err != nil {
				t.Skip("c++ compiler not found, skipping")
			}

			engine := Engines[".cpp"]
			out, err := engine.Run(testDataPath("cpp/simple.cpp"), "")
			if err != nil {
				t.Fatalf("unexpected error: %v\noutput: %s", err, out)
			}
			if out != expected {
				t.Fatalf("expected %q, got %q", expected, out)
			}
		})

		t.Run("py", func(t *testing.T) {
			if _, err := exec.LookPath("python3"); err != nil {
				t.Skip("python3 not found, skipping")
			}

			engine := Engines[".py"]
			out, err := engine.Run(testDataPath("python/simple.py"), "")
			if err != nil {
				t.Fatalf("unexpected error: %v\noutput: %s", err, out)
			}
			if out != expected {
				t.Fatalf("expected %q, got %q", expected, out)
			}
		})
	})
}
