package kubeClient

import (
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// GetKubeClient(clusterContext string) will return a kube client based
// on if you are running cli on a personal dev machine or
// within a kubernetes cluster. You can also provide a context
// to direct request a specific cluster
func GetKubeClient(ops ...options) (*kubernetes.Clientset, error) {
	clusterContext := getKubeCtx(ops...)
	var client *kubernetes.Clientset
	var err error
	switch {
	case clusterContext != "":
		client, err = getKubeClientWithContext(ops...)
	case serverRunningInCluster():
		client, err = getInternalClusterKubeClient(ops...)
	default:
		client, err = getKubeClient(ops...)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get kube cfg: %w", err)
	}

	return client, nil
}

func getInternalClusterKubeClient(ops ...options) (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, wrapErr(err, "error loading in-cluster config")
	}

	config.UserAgent = getUserAgent(ops...)
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, wrapErr(err, "error creating clientset from in-cluster config")
	}

	return clientSet, nil
}

func getKubeClientWithContext(ops ...options) (*kubernetes.Clientset, error) {
	kubeConfig, err := getKubeConfigPath()
	if err != nil {
		return nil, err
	}
	configLoadingRules := &clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeConfig}
	configOverrides := &clientcmd.ConfigOverrides{CurrentContext: getKubeCtx(ops...)}

	kConfig, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(configLoadingRules, configOverrides).ClientConfig()
	if err != nil {
		return nil, wrapErr(err, "error loading kube config. Do you have access to this cluster?")
	}

	kConfig.UserAgent = getUserAgent(ops...)

	clientSet, err := kubernetes.NewForConfig(kConfig)
	if err != nil {
		return nil, wrapErr(err, "error loading kube config. Do you have access to this cluster?")
	}

	return clientSet, nil

}

func getKubeClient(ops ...options) (*kubernetes.Clientset, error) {
	kubeConfig, err := getKubeConfigPath()
	if err != nil {
		return nil, err
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
	if err != nil {
		return nil, wrapErr(err, "error loading kube config")
	}

	config.UserAgent = getUserAgent(ops...)
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, wrapErr(err, "error loading kube config")
	}

	return clientSet, nil
}
