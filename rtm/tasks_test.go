package rtm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTasksGetList(t *testing.T) {
	ts := prepareTestServer(map[string]string{})
	defer ts.Close()
	client := prepareClient(ts)

	tasks, err := client.Tasks.GetList("", "", "")

	assert.Equal(t, nil, err)
	assert.Equal(t, 2, len(tasks))
	assert.Equal(t, "task1", tasks[0].Name)
	assert.Equal(t, "task2", tasks[1].Name)
}

func TestTasksAdd(t *testing.T) {
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

func TestTasksAddError(t *testing.T) {
	ts := prepareTestServer(map[string]string{"rtm.tasks.add": "error"})
	defer ts.Close()
	client := prepareClient(ts)

	timeline, _ := client.Timelines.Create()
	tasks, transaction, err := client.Tasks.Add(timeline, "newTaskName")

	assert.Equal(t, 0, len(tasks))
	assert.Equal(t, "", transaction.ID)
	assert.Equal(t, "Task name provided is invalid. (code=4000)", err.Error())
}

func TestTasksDelete(t *testing.T) {
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

func TestTasksDeleteError(t *testing.T) {
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

func TestTasksComplete(t *testing.T) {
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

func TestTasksUncomplete(t *testing.T) {
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

func TestTasksDueDate(t *testing.T) {
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

func TestTasksAddTags(t *testing.T) {
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
