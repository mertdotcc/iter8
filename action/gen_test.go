package action

import (
	"os"
	"testing"

	"github.com/iter8-tools/iter8/base"
	"github.com/stretchr/testify/assert"
)

func TestGen(t *testing.T) {
	// fix gOpts
	os.Chdir(t.TempDir())
	gOpts := NewGenOpts()
	gOpts.ChartsParentDir = base.CompletePath("../", "")
	gOpts.ChartName = "load-test-http"
	gOpts.Values = []string{"url=https://httpbin.org/get"}
	err := gOpts.LocalRun()
	assert.NoError(t, err)
}
