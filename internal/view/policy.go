package view

import (
	"context"

	"github.com/open-infra/osc/internal"
	"github.com/open-infra/osc/internal/client"
	"github.com/open-infra/osc/internal/render"
	"github.com/open-infra/osc/internal/ui"
	"github.com/gdamore/tcell/v2"
)

const (
	group = "Group"
	user  = "User"
	sa    = "ServiceAccount"
)

// Policy presents a RBAC rules viewer based on what a given user/group or sa can do.
type Policy struct {
	ResourceViewer

	subjectKind, subjectName string
}

// NewPolicy returns a new viewer.
func NewPolicy(app *App, subject, name string) *Policy {
	p := Policy{
		ResourceViewer: NewBrowser(client.NewGVR("policy")),
		subjectKind:    subject,
		subjectName:    name,
	}
	p.GetTable().SetColorerFn(render.Policy{}.ColorerFunc())
	p.AddBindKeysFn(p.bindKeys)
	p.GetTable().SetSortCol(nameCol, false)
	p.SetContextFn(p.subjectCtx)
	p.GetTable().SetEnterFn(blankEnterFn)

	return &p
}

func (p *Policy) subjectCtx(ctx context.Context) context.Context {
	ctx = context.WithValue(ctx, internal.KeySubjectKind, mapSubject(p.subjectKind))
	ctx = context.WithValue(ctx, internal.KeyPath, mapSubject(p.subjectKind)+":"+p.subjectName)
	return context.WithValue(ctx, internal.KeySubjectName, p.subjectName)
}

func (p *Policy) bindKeys(aa ui.KeyActions) {
	aa.Delete(ui.KeyShiftA, tcell.KeyCtrlSpace, ui.KeySpace)
	aa.Add(ui.KeyActions{
		ui.KeyShiftN: ui.NewKeyAction("Sort Name", p.GetTable().SortColCmd(nameCol, true), false),
		ui.KeyShiftO: ui.NewKeyAction("Sort Group", p.GetTable().SortColCmd("GROUP", true), false),
		ui.KeyShiftB: ui.NewKeyAction("Sort Binding", p.GetTable().SortColCmd("BINDING", true), false),
	})
}

func mapSubject(subject string) string {
	switch subject {
	case "g":
		return group
	case "s":
		return sa
	case "u":
		return user
	default:
		return subject
	}
}
