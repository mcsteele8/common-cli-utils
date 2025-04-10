package examples

import "github.com/mcsteele8/common-cli-utils/kubeClient"

func exampleKubeClient() {
	kubeClient.GetKubeClient(kubeClient.SetUserAgent("test_user_agent"))
}
