package kube

import (
	"github.com/pkg/errors"
	"github.com/wtfutil/wtf/wtf"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type KubeInfo struct {
	Namespace    string
	Context      string
	Config       *restclient.Config
	ClientSet    *kubernetes.Clientset
	Pods         *v1.PodList
	Nodes        *v1.NodeList
	Error        bool
	ErrorMessage string
}

func NewKube() *KubeInfo {

	clientConfig := clientcmd.NewDefaultClientConfigLoadingRules()

	// use the current context in kubeconfig
	configPath := wtf.Config.UString("wtf.mods.kube.configpath", "")
	if configPath != "" {
		clientConfig.ExplicitPath = configPath
	}

	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(clientConfig,
		&clientcmd.ConfigOverrides{CurrentContext: ""}).ClientConfig()

	clientset, err := kubernetes.NewForConfig(config)

	if err != nil {
		return &KubeInfo {
			Error: true,
			ErrorMessage: "Unable to load kube config",
		}
	}

	namespace := wtf.Config.UString("wtf.mods.kube.namespace", "")

	repo := KubeInfo{
		Namespace: namespace,
		Config:    config,
		ClientSet: clientset,
		Error: false,
	}
	return &repo
}

// Refresh updates pods and nodes
func (repo *KubeInfo) Refresh() {
	pods, err := repo.loadPods()
	if err != nil {
		repo.Pods = nil
	}
	repo.Pods = pods
	nodes, err := repo.loadNodes()
	if err != nil {
		repo.Nodes = nil
	}
	repo.Nodes = nodes
	return
}

func (repo *KubeInfo) NodeCount() int {
	if repo.Nodes == nil {
		return 0
	}
	return len(repo.Nodes.Items)
}

func (repo *KubeInfo) PodCount() int {
	if repo.Pods == nil {
		return 0
	}
	return len(repo.Pods.Items)
}

func (repo *KubeInfo) HealthyPods() (*v1.PodList, error) {
	return repo.podInfo(true)
}

func (repo *KubeInfo) UnHealthyPods() (*v1.PodList, error) {
	return repo.podInfo(false)
}

func (repo *KubeInfo) podInfo(healthy bool) (*v1.PodList, error) {
	if repo.Pods == nil {
		return nil, errors.New("No pods error")
	}
	podList := &v1.PodList{}
	for _, p := range repo.Pods.Items {
		if healthy {
			if repo.isHealthy(p) {
				podList.Items = append(podList.Items, p)
			}
		} else {
			if !repo.isHealthy(p) {
				podList.Items = append(podList.Items, p)
			}
		}
	}
	return podList, nil
}

// This may need tuning
func (repo *KubeInfo) isHealthy(pod v1.Pod) bool {
	if pod.Status.ContainerStatuses != nil && len(pod.Status.ContainerStatuses) > 0 {
		if pod.Status.ContainerStatuses[0].State.Waiting == nil &&
			pod.Status.Phase == v1.PodRunning {
			return true
		}
	}
	return false
}

/* -------------------- Unexported Functions -------------------- */

func (repo *KubeInfo) loadPods() (*v1.PodList, error) {
	if repo == nil || repo.ClientSet == nil {
		return nil, errors.New("Unable to connect to kube")
	}
	pods, err := repo.ClientSet.CoreV1().Pods(repo.Namespace).List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return pods, nil
}

func (repo *KubeInfo) loadNodes() (*v1.NodeList, error) {
	nodes, err := repo.ClientSet.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return nodes, nil
}

