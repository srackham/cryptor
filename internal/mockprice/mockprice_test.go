package mockprice

import (
	"testing"

	"github.com/srackham/cryptor/internal/assert"
	"github.com/srackham/cryptor/internal/helpers"
)

func (r *Reader) TestPrice(t *testing.T) {
	// TODO test Symbol XXX returns error
	reader := Reader{}
	amt, _ := reader.GetPrice("BTC", helpers.DateNowString())
	assert.Equal(t, 10000, amt)
}
