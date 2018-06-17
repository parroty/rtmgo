package rtm

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func prepareTestServer(modes map[string]string) *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc(
		"/services/rest",
		func(w http.ResponseWriter, r *http.Request) {
			q := r.URL.Query()
			method := q["method"]
			//fmt.Println(r.URL)
			//fmt.Println()

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

func TestGetList(t *testing.T) {
	ts := prepareTestServer(map[string]string{})
	defer ts.Close()
	client := prepareClient(ts)

	tasks, err := client.Tasks.GetList("", "", "")

	assert.Equal(t, nil, err)
	assert.Equal(t, 2, len(tasks))
	assert.Equal(t, "task1", tasks[0].Name)
	assert.Equal(t, "task2", tasks[1].Name)
}

func TestAdd(t *testing.T) {
	ts := prepareTestServer(map[string]string{})
	defer ts.Close()
	client := prepareClient(ts)

	timeline, _ := client.Timelines.Create()
	tasks, transaction, err := client.Tasks.Add(timeline, "newTaskName")

	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(tasks))
	assert.Equal(t, "newTaskName", tasks[0].Name)
	assert.Equal(t, "200000000", transaction.ID)
}

func TestAddError(t *testing.T) {
	ts := prepareTestServer(map[string]string{"rtm.tasks.add": "error"})
	defer ts.Close()
	client := prepareClient(ts)

	timeline, _ := client.Timelines.Create()
	tasks, transaction, err := client.Tasks.Add(timeline, "newTaskName")

	assert.Equal(t, 0, len(tasks))
	assert.Equal(t, "", transaction.ID)
	assert.Equal(t, "Task name provided is invalid. (code=4000)", err.Error())
}

func TestDelete(t *testing.T) {
	ts := prepareTestServer(map[string]string{})
	defer ts.Close()
	client := prepareClient(ts)

	chunk := Chunk{}
	task := Task{Chunks: []Chunk{chunk}}
	timeline, _ := client.Timelines.Create()
	transactions, err := client.Tasks.Delete(task, timeline)

	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(transactions))
}

func TestDeleteError(t *testing.T) {
	ts := prepareTestServer(map[string]string{"rtm.tasks.delete": "error"})
	defer ts.Close()
	client := prepareClient(ts)

	chunk := Chunk{}
	task := Task{Chunks: []Chunk{chunk}}
	timeline, _ := client.Timelines.Create()
	transactions, err := client.Tasks.Delete(task, timeline)

	assert.Equal(t, 0, len(transactions))
	assert.Equal(t, "taskseries_id/task_id invalid or not provided (code=340)", err.Error())
}

func TestComplete(t *testing.T) {
	ts := prepareTestServer(map[string]string{})
	defer ts.Close()
	client := prepareClient(ts)

	chunk := Chunk{}
	task := Task{Chunks: []Chunk{chunk}}
	timeline, _ := client.Timelines.Create()
	transactions, err := client.Tasks.Complete(task, timeline)

	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(transactions))
}

func TestUncomplete(t *testing.T) {
	ts := prepareTestServer(map[string]string{})
	defer ts.Close()
	client := prepareClient(ts)

	chunk := Chunk{}
	task := Task{Chunks: []Chunk{chunk}}
	timeline, _ := client.Timelines.Create()
	transactions, err := client.Tasks.Uncomplete(task, timeline)

	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(transactions))
}

func TestDueDate(t *testing.T) {
	ts := prepareTestServer(map[string]string{})
	defer ts.Close()
	client := prepareClient(ts)

	chunk := Chunk{}
	task := Task{Chunks: []Chunk{chunk}}
	timeline, _ := client.Timelines.Create()
	transactions, err := client.Tasks.SetDueDate(task, timeline, "tomorrow")

	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(transactions))
}

func TestAddTags(t *testing.T) {
	ts := prepareTestServer(map[string]string{})
	defer ts.Close()
	client := prepareClient(ts)

	chunk := Chunk{}
	task := Task{Chunks: []Chunk{chunk}}
	timeline, _ := client.Timelines.Create()
	transactions, err := client.Tasks.AddTags(task, timeline, "a,b,c")

	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(transactions))
}
