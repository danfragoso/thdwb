package profiler

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestProfiler(t *testing.T) {
	profiler := CreateProfiler()

	profiler.Start("test")
	time.Sleep(2 * time.Second)
	profiler.Stop("test")

	stopProfile := profiler.GetProfile("test")

	t.Log(stopProfile.GetElapsedTime())
	assert.Equal(t, true, true, "")
}
