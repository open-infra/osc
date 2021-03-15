package view_test

import (
	"context"
	"testing"

	"github.com/open-infra/osc/internal"
	"github.com/open-infra/osc/internal/client"
	"github.com/open-infra/osc/internal/config"
	"github.com/open-infra/osc/internal/view"
	"github.com/stretchr/testify/assert"
)

func TestPodNew(t *testing.T) {
	po := view.NewPod(client.NewGVR("v1/pods"))

	assert.Nil(t, po.Init(makeCtx()))
	assert.Equal(t, "Pods", po.Name())
	assert.Equal(t, 24, len(po.Hints()))
}

// Helpers...

func makeCtx() context.Context {
	cfg := config.NewConfig(ks{})
	return context.WithValue(context.Background(), internal.KeyApp, view.NewApp(cfg))
}
