package kubeClient

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

const (
	kubeContext   kubeClientOpsTypes = 0
	kubeUserAgent kubeClientOpsTypes = 1
)

type kubeClientOpsTypes int

type options struct {
	id    kubeClientOpsTypes
	value string
}

func SetContext(ctx string) options {
	if ctx == "" {
		return options{
			value: "",
			id:    kubeContext,
		}
	}
	return options{
		value: ctx,
		id:    kubeContext,
	}
}

func SetUserAgent(userAgent string) options {
	if userAgent == "" {
		return options{
			value: "",
			id:    kubeUserAgent,
		}
	}
	return options{
		value: userAgent,
		id:    kubeUserAgent,
	}
}

func getKubeCtx(ops ...options) string {
	for _, o := range ops {
		if o.id == kubeContext {
			return o.value
		}
	}
	return ""
}

func getUserAgent(ops ...options) string {
	for _, o := range ops {
		if o.id == kubeUserAgent {
			return o.value
		}
	}
	return ""
}

func serverRunningInCluster() bool {
	// env var gets injected on all new pods
	// https://kubernetes.io/docs/tutorials/services/connect-applications-service/#environment-variables
	return os.Getenv("KUBERNETES_SERVICE_HOST") != ""
}

func getKubeConfigPath() (string, error) {
	var kubeConfigPath *string
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln("failed to get user home directory: ", err.Error())
	}
	switch {
	case os.Getenv("KUBECONFIG") != "":
		kubeConfigPath = flag.String("kubeconfig", os.Getenv("KUBECONFIG"), "(optional) absolute path to the kubeconfig file")
	case homeDir != "":
		kubeConfigPath = flag.String("kubeconfig", filepath.Join(homeDir, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	default:
		kubeConfigPath = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}

	return *kubeConfigPath, nil
}

func wrapErr(err error, msg string) error {
	return fmt.Errorf("%s: %w", msg, err)
}
