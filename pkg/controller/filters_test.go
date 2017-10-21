package controller

import (
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestFilterPodsEmpty(t *testing.T) {
	assert := assert.New(t)

	reducedPods := TestRunFilter([]*v1.Pod{}, "tr1")
	assert.Equal(len(reducedPods), 0, "The length after filtering should be 0")
}

func TestFilterPods(t *testing.T) {
	assert := assert.New(t)
	pods := []*v1.Pod{
		&v1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Labels: map[string]string{"test-run": "tr1"},
			},
		},
		&v1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Labels: map[string]string{"test-run": "tr2"},
			},
		},
	}
	reducedPods := TestRunFilter(pods, "tr1")
	assert.Equal(len(reducedPods), 1, "Should only have one test with testrun tr1")
}
