package kubeClient

import (
	istioClient "istio.io/client-go/pkg/clientset/versioned"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Istio kube clients are used to interact with crd. For example we need this
// to be able to access argo CRDs because the sdk for argo has a bunch of bad dependancies in it
// so we have to use
func GetIstioKubeClient(ops ...options) (*istioClient.Clientset, error) {
	clusterContext := getKubeCtx(ops...)
	var client *istioClient.Clientset
	var err error
	switch {
	case clusterContext != "":
		client, err = getIstioKubeClientWithContext(ops...)
	case serverRunningInCluster():
		client, err = getInternalClusterIstioKubeClient(ops...)
	default:
		client, err = getIstioKubeClient(ops...)
	}

	if err != nil {
		return nil, wrapErr(err, "failed to get Istio kube config")
	}

	return client, nil
}

func getInternalClusterIstioKubeClient(ops ...options) (*istioClient.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, wrapErr(err, "error loading in-cluster config")
	}

	config.UserAgent = getUserAgent(ops...)
	client, err := istioClient.NewForConfig(config)
	if err != nil {
		return nil, wrapErr(err, "error loading kube config")
	}

	return client, nil
}

func getIstioKubeClientWithContext(ops ...options) (*istioClient.Clientset, error) {
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
	client, err := istioClient.NewForConfig(config)
	if err != nil {
		return nil, wrapErr(err, "error loading kube config")
	}

	return client, nil

}

func getIstioKubeClient(ops ...options) (*istioClient.Clientset, error) {
	kubeconfig, err := getKubeConfigPath()
	if err != nil {
		return nil, err
	}
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, wrapErr(err, "error loading kube config")
	}

	config.UserAgent = getUserAgent(ops...)
	client, err := istioClient.NewForConfig(config)
	if err != nil {
		return nil, wrapErr(err, "error loading kube config")
	}

	return client, err
}
