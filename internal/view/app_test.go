package view_test

import (
	"testing"

	"github.com/open-infra/osc/internal/config"
	"github.com/open-infra/osc/internal/view"
	"github.com/stretchr/testify/assert"
)

func TestAppNew(t *testing.T) {
	a := view.NewApp(config.NewConfig(ks{}))
	a.Init("blee", 10)

	assert.Equal(t, 10, len(a.GetActions()))
}
