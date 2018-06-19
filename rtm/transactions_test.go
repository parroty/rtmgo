package rtm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransactionsUndo(t *testing.T) {
	ts := prepareTestServer(map[string]string{})
	defer ts.Close()
	client := prepareClient(ts)

	transaction := Transaction{}
	timeline, _ := client.Timelines.Create()
	err := client.Transactions.Undo(timeline, transaction)

	assert.Equal(t, nil, err)
}
