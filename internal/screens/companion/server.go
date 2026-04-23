package companion

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
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

type CompanionData struct {
	Name  string `json:"name"`
	Group string `json:"group"`
	URL   string `json:"url"`
}

func CreateServer(onData func(CompanionData)) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		defer r.Body.Close()

		var data CompanionData
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
