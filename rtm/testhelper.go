package rtm

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

func prepareTestServer(modes map[string]string) *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc(
		"/services/rest",
		func(w http.ResponseWriter, r *http.Request) {
			q := r.URL.Query()
			method := q["method"]

			var name = method[0]
			if val, ok := modes[name]; ok {
				name = name + "." + val
			}
			path := "../testdata/" + name + ".json"
			data, err := ioutil.ReadFile(path)
			if err != nil {
				panic(err)
			}

			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, string(data))
		},
	)
	return httptest.NewServer(mux)
}

func prepareClient(ts *httptest.Server) *Client {
	client := NewClient("dummyToken", "dummyApiKey")
	client.SetBaseURL(ts.URL)
	return client
}
