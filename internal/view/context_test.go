package view_test

import (
	"testing"

	"github.com/open-infra/osc/internal/client"
	"github.com/open-infra/osc/internal/view"
	"github.com/stretchr/testify/assert"
)

func TestContext(t *testing.T) {
	ctx := view.NewContext(client.NewGVR("contexts"))

	assert.Nil(t, ctx.Init(makeCtx()))
	assert.Equal(t, "Contexts", ctx.Name())
	assert.Equal(t, 4, len(ctx.Hints()))
}
