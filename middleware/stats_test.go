package middleware

import (
	"testing"

	"net/http"

	"github.com/hellofresh/janus/pkg/test"
	stats "github.com/hellofresh/stats-go"
	"github.com/stretchr/testify/assert"
)

func TestSuccessfulStats(t *testing.T) {
	mw := NewStats(stats.NewStatsdClient("", ""))
	w, err := test.Record(
		"GET",
		"/",
		map[string]string{
			"Content-Type": "application/json",
		},
		recovery(mw.Handler(http.HandlerFunc(test.Ping))),
	)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
}
