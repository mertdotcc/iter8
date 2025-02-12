package action

import (
	"context"
	"os"
	"testing"

	"github.com/iter8-tools/iter8/base"
	"github.com/iter8-tools/iter8/driver"
	"github.com/stretchr/testify/assert"
	"helm.sh/helm/v3/pkg/cli"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestKubeAssert(t *testing.T) {
	_ = os.Chdir(t.TempDir())
	// fix aOpts
	aOpts := NewAssertOpts(driver.NewFakeKubeDriver(cli.New()))
	aOpts.Conditions = []string{Completed, NoFailure, SLOs}

	byteArray, _ := os.ReadFile(base.CompletePath("../testdata/assertinputs", driver.ExperimentPath))
	_, _ = aOpts.Clientset.CoreV1().Secrets("default").Create(context.TODO(), &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "default",
			Namespace: "default",
		},
		StringData: map[string]string{driver.ExperimentPath: string(byteArray)},
	}, metav1.CreateOptions{})

	ok, err := aOpts.KubeRun()
	assert.True(t, ok)
	assert.NoError(t, err)
}
