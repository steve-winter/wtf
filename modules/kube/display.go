package kube

import (
	"fmt"

	"k8s.io/api/core/v1"
)

func (widget *Widget) display() {
	kube := widget.currentKubeInfo()
	if kube == nil {
		widget.View.SetText(" [red]Unable to retrieve kube info[white]")
		return
	}

	if kube.Error {
		widget.View.SetText(fmt.Sprintf(" [red]%s[white]", kube.ErrorMessage))
		return
	}

	widget.View.SetTitle(widget.ContextualTitle(fmt.Sprintf("%s - %s", widget.Name, widget.title(kube))))
	if kube.Nodes == nil && kube.Pods == nil {
		widget.View.SetText(" [red]Unable to retrieve kube info[white]")
	} else {
		str := widget.displayPodCount(kube)
		str = str + widget.displayHealthyPodInfo(kube)
		str = str + widget.displayUnHealthyPodInfo(kube)
		str = str + widget.displayNodeCount(kube)
		widget.View.SetText(str)
	}
}

func (widget *Widget) displayNodeCount(info *KubeInfo) string {
	str := ""
	str = str + fmt.Sprintf(" Node Count:[yellow] %d[white]\n", info.NodeCount())
	return str
}

func (widget *Widget) displayPods(info *KubeInfo) string {
	str := ""
	for _, p := range info.Pods.Items {
		str = str + fmt.Sprintf("[blue]%s[white]\n", p.Name)

	}
	return str
}

func (widget *Widget) displayHealthyPodInfo(info *KubeInfo) string {
	healthy, err := info.HealthyPods()
	if err != nil {
		return ""
	}
	return widget.podStatusView(healthy, "Healthy")
}

func (widget *Widget) displayUnHealthyPodInfo(info *KubeInfo) string {
	healthy, err := info.UnHealthyPods()
	if err != nil {
		return ""
	}
	return widget.podStatusView(healthy, "UnHealthy")
}

func (widget *Widget) podStatusView(pods *v1.PodList, status string) string {
	str := ""
	col := "[red]"
	if status == "Healthy" {
		col = "[green]"
	}
	str = str + fmt.Sprintf("[white] %s Pods: %s%d [white]\n", status, col, len(pods.Items))
	return str
}

func (widget *Widget) displayPodCount(info *KubeInfo) string {
	str := ""
	str = str + fmt.Sprintf(" Pod Count:[yellow] %d[white]\n", info.PodCount())
	return str
}

func (widget *Widget) title(repo *KubeInfo) string {
	title := repo.Namespace
	if title == "" {
		title = "ALL"
	}
	return fmt.Sprintf("[green]%s[white]", title)
}

