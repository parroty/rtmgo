package rtm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListsGetList(t *testing.T) {
	ts := prepareTestServer(map[string]string{})
	defer ts.Close()
	client := prepareClient(ts)

	tasks, err := client.Lists.GetList()

	assert.Equal(t, nil, err)
	assert.Equal(t, 3, len(tasks))
	assert.Equal(t, "Inbox", tasks[0].Name)
	assert.Equal(t, "All Tasks", tasks[1].Name)
	assert.Equal(t, "newListName", tasks[2].Name)
}

func TestListsAdd(t *testing.T) {
	ts := prepareTestServer(map[string]string{})
	defer ts.Close()
	client := prepareClient(ts)

	timeline, _ := client.Timelines.Create()
	list, transaction, err := client.Lists.Add(timeline, "newListName")

	assert.Equal(t, nil, err)
	assert.Equal(t, "newListName", list.Name)
	assert.Equal(t, "600000000", transaction.ID)
}

func TestListsDelete(t *testing.T) {
	ts := prepareTestServer(map[string]string{})
	defer ts.Close()
	client := prepareClient(ts)

	list := List{Name: "newListName"}
	timeline, _ := client.Timelines.Create()
	deletedList, transaction, err := client.Lists.Delete(list, timeline)

	assert.Equal(t, nil, err)
	assert.Equal(t, "newListName", deletedList.Name)
	assert.Equal(t, "1", deletedList.Deleted)
	assert.Equal(t, "600000000", transaction.ID)
}
