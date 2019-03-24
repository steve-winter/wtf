package kube

import (
	"github.com/rivo/tview"
	"github.com/wtfutil/wtf/wtf"
)

const HelpText = `
 Keyboard commands for Kube:

   Retrieves Pod and Node info using Kubectl
`

type Widget struct {
	wtf.HelpfulWidget
	wtf.TextWidget

	KubeInfos []*KubeInfo
	Idx         int
	UpdateCount int
}

func NewWidget(app *tview.Application, pages *tview.Pages) *Widget {
	widget := Widget{
		HelpfulWidget: wtf.NewHelpfulWidget(app, pages, HelpText),
		TextWidget:    wtf.NewTextWidget(app,"Kube", "kube", true),
		Idx: 0,
		UpdateCount: 0,
	}

	widget.KubeInfos = widget.buildKubeCollection()
	widget.TextWidget.RefreshInt = wtf.Config.UInt("wtf.mods.kube.refreshtime", 60)
	widget.HelpfulWidget.SetView(widget.View)
	widget.View.SetText("Loading...")
	widget.View.SetRegions(true)

	return &widget
}

/* -------------------- Exported Functions -------------------- */

func (widget *Widget) Refresh() {
	for _, repo := range widget.KubeInfos {
		repo.Refresh()
	}
	widget.UpdateCount = widget.UpdateCount+1
	widget.display()
}

func (widget *Widget) Next() {
	widget.Idx = widget.Idx + 1
	if widget.Idx == len(widget.KubeInfos) {
		widget.Idx = 0
	}

	widget.display()
}

func (widget *Widget) Prev() {
	widget.Idx = widget.Idx - 1
	if widget.Idx < 0 {
		widget.Idx = len(widget.KubeInfos) - 1
	}

	widget.display()
}

/* -------------------- Unexported Functions -------------------- */

//Maintained as array to allow future Namespace/Context cycling
func (widget *Widget) buildKubeCollection() []*KubeInfo {
	KubeInfos := []*KubeInfo{}

	repo := NewKube()
	KubeInfos = append(KubeInfos, repo)
	return KubeInfos
}

func (widget *Widget) currentKubeInfo() *KubeInfo {
	if len(widget.KubeInfos) == 0 {
		return nil
	}
	if widget.Idx < 0 || widget.Idx >= len(widget.KubeInfos) {
		return nil
	}
	return widget.KubeInfos[widget.Idx]
}

