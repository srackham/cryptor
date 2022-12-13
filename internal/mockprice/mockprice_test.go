package mockprice

import (
	"testing"

	"github.com/srackham/cryptor/internal/assert"
)

func (r *Reader) TestPrice(t *testing.T) {
	// TODO test Symbol XXX returns error
	reader := Reader{}
	amt, _ := reader.GetPrice("BTC", "latest")
	assert.Equal(t, 10000, amt)
}
