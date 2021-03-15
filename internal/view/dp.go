package view

import (
	"github.com/open-infra/osc/internal/client"
	"github.com/open-infra/osc/internal/dao"
	"github.com/open-infra/osc/internal/render"
	"github.com/open-infra/osc/internal/ui"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const scaleDialogKey = "scale"

// Deploy represents a deployment view.
type Deploy struct {
	ResourceViewer
}

// NewDeploy returns a new deployment view.
func NewDeploy(gvr client.GVR) ResourceViewer {
	d := Deploy{
		ResourceViewer: NewPortForwardExtender(
			NewRestartExtender(
				NewScaleExtender(
					NewImageExtender(
						NewLogsExtender(
							NewBrowser(gvr),
							nil,
						),
					),
				),
			),
		),
	}
	d.AddBindKeysFn(d.bindKeys)
	d.GetTable().SetEnterFn(d.showPods)
	d.GetTable().SetColorerFn(render.Deployment{}.ColorerFunc())

	return &d
}

func (d *Deploy) bindKeys(aa ui.KeyActions) {
	aa.Add(ui.KeyActions{
		ui.KeyShiftR: ui.NewKeyAction("Sort Ready", d.GetTable().SortColCmd(readyCol, true), false),
		ui.KeyShiftU: ui.NewKeyAction("Sort UpToDate", d.GetTable().SortColCmd(uptodateCol, true), false),
		ui.KeyShiftL: ui.NewKeyAction("Sort Available", d.GetTable().SortColCmd(availCol, true), false),
	})
}

func (d *Deploy) showPods(app *App, model ui.Tabular, gvr, path string) {
	var ddp dao.Deployment
	dp, err := ddp.Load(app.factory, path)
	if err != nil {
		app.Flash().Err(err)
		return
	}

	showPodsFromSelector(app, path, dp.Spec.Selector)
}

// ----------------------------------------------------------------------------
// Helpers...

func showPodsFromSelector(app *App, path string, sel *metav1.LabelSelector) {
	l, err := metav1.LabelSelectorAsSelector(sel)
	if err != nil {
		app.Flash().Err(err)
		return
	}

	showPods(app, path, l.String(), "")
}
