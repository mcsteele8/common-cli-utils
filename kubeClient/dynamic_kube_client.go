package kubeClient

import (
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Dynamic kube clients are used to interact with crd. For example we need this
// to be able to access argo CRDs because the sdk for argo has a bunch of bad dependencies in it
// so we have to use
func GetDynamicKubeClient(ops ...options) (*dynamic.DynamicClient, error) {
	clusterContext := getKubeCtx(ops...)
	var client *dynamic.DynamicClient
	var err error
	switch {
	case clusterContext != "":
		client, err = getDynamicKubeClientWithContext(ops...)
	case serverRunningInCluster():
		client, err = getInternalClusterDynamicKubeClient(ops...)
	default:
		client, err = getDynamicKubeClient(ops...)
	}

	if err != nil {
		return nil, wrapErr(err, "failed to get dynamic kube config")
	}

	return client, nil
}

func getInternalClusterDynamicKubeClient(ops ...options) (*dynamic.DynamicClient, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, wrapErr(err, "error loading in-cluster config")
	}

	config.UserAgent = getUserAgent(ops...)
	client, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, wrapErr(err, "error loading kube config")
	}

	return client, nil
}

func getDynamicKubeClientWithContext(ops ...options) (*dynamic.DynamicClient, error) {
	kubeconfig, err := getKubeConfigPath()
	if err != nil {
		return nil, err
	}

	configLoadingRules := &clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfig}
	configOverrides := &clientcmd.ConfigOverrides{CurrentContext: getKubeCtx(ops...)}

	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(configLoadingRules, configOverrides).ClientConfig()
	if err != nil {
		return nil, wrapErr(err, "error loading kube config. Do you have access to this cluster?")
	}

	config.UserAgent = getUserAgent(ops...)
	client, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, wrapErr(err, "error loading kube config")
	}

	return client, nil

}

func getDynamicKubeClient(ops ...options) (*dynamic.DynamicClient, error) {
	kubeconfig, err := getKubeConfigPath()
	if err != nil {
		return nil, err
	}
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, wrapErr(err, "error loading kube config")
	}

	config.UserAgent = getUserAgent(ops...)
	client, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, wrapErr(err, "error loading kube config")
	}

	return client, err
}
