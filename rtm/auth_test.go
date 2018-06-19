package rtm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthGetFrob(t *testing.T) {
	ts := prepareTestServer(map[string]string{})
	defer ts.Close()
	client := prepareClient(ts)

	frob, err := client.Auth.GetFrob()

	assert.Equal(t, nil, err)
	assert.Equal(t, "a267489abc8806a623e65c26cce48b8f60795ef2", frob)
}
