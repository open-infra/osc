package view_test

import (
	"testing"

	"github.com/open-infra/osc/internal/client"
	"github.com/open-infra/osc/internal/view"
	"github.com/stretchr/testify/assert"
)

func TestNSCleanser(t *testing.T) {
	ns := view.NewNamespace(client.NewGVR("v1/namespaces"))

	assert.Nil(t, ns.Init(makeCtx()))
	assert.Equal(t, "Namespaces", ns.Name())
	assert.Equal(t, 6, len(ns.Hints()))
}
