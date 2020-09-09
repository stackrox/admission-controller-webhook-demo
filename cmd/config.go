package main

import (
	"log"
	"os"
	"strings"
	"sync"

	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"
)

// Config contains runtime configuration needed for the controller.
type Config struct {
	IgnoredNamespaces         []string
	IgnoreNamespaceAnnotation string
	Namespaces                v1.NamespaceInterface
}

var (
	config     *Config
	configSync sync.Once
)

// GetConfig returns Config struct, and initializes it first time
func GetConfig() *Config {
	configSync.Do(func() {
		clientSet, err := newClientSet()
		if err != nil {
			log.Fatal(err)
		}
		config = &Config{
			IgnoredNamespaces:         strings.Split(os.Getenv("IGNORED_NAMESPACES"), ","),
			IgnoreNamespaceAnnotation: os.Getenv("IGNORE_NAMESPACE_ANNOTATION"),
			Namespaces:                clientSet.CoreV1().Namespaces(),
		}
	})

	return config
}

func newClientSet() (*kubernetes.Clientset, error) {
	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{}).ClientConfig()
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}
