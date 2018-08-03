package client

import (
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

func findExistingKubeConfig() string {
	// locations using DefaultClientConfigLoadingRules
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	return rules.GetDefaultFilename()
}

// CreateKubeClientSet obtains the client set interface from Kubeconfig
func CreateKubeClientSet() (clientset.Interface, error) {
	//TODO: It should implement only once
	//rather than once per test
	filename := findExistingKubeConfig()
	c := clientcmd.GetConfigFromFileOrDie(filename)
	restConfig, err := clientcmd.NewDefaultClientConfig(*c, &clientcmd.ConfigOverrides{ClusterInfo: clientcmdapi.Cluster{Server: ""}}).ClientConfig()
	if err != nil {
		return nil, err
	}
	clientSet, err := clientset.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}
	return clientSet, nil
}
