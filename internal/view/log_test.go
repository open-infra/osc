package view_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/derailed/tview"
	"github.com/open-infra/osc/internal/client"
	"github.com/open-infra/osc/internal/config"
	"github.com/open-infra/osc/internal/dao"
	"github.com/open-infra/osc/internal/view"
	"github.com/stretchr/testify/assert"
)

func TestLog(t *testing.T) {
	v := view.NewLog(client.NewGVR("v1/pods"), "fred/p1", "blee", false)
	v.Init(makeContext())

	v.Flush(dao.LogItems{
		dao.NewLogItemFromString("blee"),
		dao.NewLogItemFromString("bozo"),
	}.Lines(false))

	assert.Equal(t, 29, len(v.Logs().GetText(true)))
}

func BenchmarkLogFlush(b *testing.B) {
	v := view.NewLog(client.NewGVR("v1/pods"), "fred/p1", "blee", false)
	v.Init(makeContext())

	items := dao.LogItems{
		dao.NewLogItemFromString("blee"),
		dao.NewLogItemFromString("bozo"),
	}
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		v.Flush(items.Lines(false))
	}
}

func TestLogAnsi(t *testing.T) {
	buff := bytes.NewBufferString("")
	w := tview.ANSIWriter(buff, "white", "black")
	fmt.Fprintf(w, "[YELLOW] ok")
	assert.Equal(t, "[YELLOW] ok", buff.String())

	v := tview.NewTextView()
	v.SetDynamicColors(true)
	aw := tview.ANSIWriter(v, "white", "black")
	s := "[2019-03-27T15:05:15,246][INFO ][o.e.c.r.a.AllocationService] [es-0] Cluster health status changed from [YELLOW] to [GREEN] (reason: [shards started [[.monitoring-es-6-2019.03.27][0]]"
	fmt.Fprintf(aw, "%s", s)
	assert.Equal(t, s+"\n", v.GetText(false))
}

func TestLogViewSave(t *testing.T) {
	v := view.NewLog(client.NewGVR("v1/pods"), "fred/p1", "blee", false)
	v.Init(makeContext())

	app := makeApp()
	v.Flush(dao.LogItems{
		dao.NewLogItemFromString("blee"),
		dao.NewLogItemFromString("bozo"),
	}.Lines(false))
	config.OscDumpDir = "/tmp"
	dir := filepath.Join(config.OscDumpDir, app.Config.Osc.CurrentCluster)
	c1, _ := ioutil.ReadDir(dir)
	v.SaveCmd(nil)
	c2, _ := ioutil.ReadDir(dir)
	assert.Equal(t, len(c2), len(c1)+1)
}

// ----------------------------------------------------------------------------
// Helpers...

func makeApp() *view.App {
	return view.NewApp(config.NewConfig(ks{}))
}
