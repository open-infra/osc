package view_test

import (
	"testing"

	"github.com/open-infra/osc/internal/client"
	"github.com/open-infra/osc/internal/view"
	"github.com/stretchr/testify/assert"
)

func TestPortForwardNew(t *testing.T) {
	pf := view.NewPortForward(client.NewGVR("portforwards"))

	assert.Nil(t, pf.Init(makeCtx()))
	assert.Equal(t, "PortForwards", pf.Name())
	assert.Equal(t, 10, len(pf.Hints()))
}
