package sessions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetractSessionId(t *testing.T) {
	key := "/btfs/16Uiu2HAkyUxz9mhH9yRg3GnLPn9DMaY1A8Ce23jAGyN6LX2XgRmz/renter/sessions/0fb2f98b-3ff2-42ca-b297-7e5e13d0fe5a/status"
	id := getSessionId(key)
	assert.Equal(t, "0fb2f98b-3ff2-42ca-b297-7e5e13d0fe5a", id)
}
