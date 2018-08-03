package app

import (
	"github.com/fseldow/AKS_daily_cleanup/pkg/client"
	clientset "k8s.io/client-go/kubernetes"
)

// Framework is the client context
type Framework struct {
	Kubeclientset clientset.Interface
	Azureclient   *client.AzureTestClient
}

func (f *Framework) getClient() (err error) {
	f.Kubeclientset, err = client.CreateKubeClientSet()
	f.Azureclient, err = client.CreateAzureTestClient()
	return
}

// NewFramework returns new default framework
func NewFramework() (f *Framework) {
	f = &Framework{}
	f.getClient()
	return
}
