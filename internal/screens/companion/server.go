package companion

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"

	"github.com/thecomputerm/lazycph/internal/core"
)

var companionPorts = []int{
	1327,
	4244,
	6174,
	10042,
	10043,
	10045,
	27121,
}

type Data struct {
	Name  string `json:"name"`
	Group string `json:"group"`
	URL   string `json:"url"`

	Tests []struct {
		Input  string `json:"input"`
		Output string `json:"output"`
	} `json:"tests"`
}

func (d Data) toTestCaseList() core.TestCaseList {
	testCases := make(core.TestCaseList, 0, len(d.Tests))
	for _, test := range d.Tests {
		tc := core.NewTestCase()
		tc.Input = test.Input
		tc.Expected = test.Output
		testCases = append(testCases, tc)
	}
	return testCases
}

func CreateServer(onData func(Data)) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		defer r.Body.Close()

		var data Data
		if err := json.NewDecoder(io.LimitReader(r.Body, 1<<20)).Decode(&data); err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		if onData != nil {
			onData(data)
		}

		w.WriteHeader(http.StatusNoContent)
	})

	for _, port := range companionPorts {
		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err != nil {
			continue
		}

		log.Fatal(http.Serve(listener, handler))
	}
}

