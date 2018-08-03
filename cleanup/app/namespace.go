package app

import (
	"github.com/fseldow/AKS_daily_cleanup/pkg/utils"
)

const (
	testNamespacePrefix = "e2e-tests"
)

// CleanupNamespace clean up all test namespace
func (f *Framework) CleanupNamespace() error {
	utils.Logf("Cleaning up azure test namespace")
	namespaceList, err := utils.GetNamespaceList(f.Kubeclientset)
	for _, namespace := range namespaceList.Items {
		if isTestNamespace(namespace.Name) {
			utils.DeleteNameSpace(f.Kubeclientset, namespace.Name)
		}
	}
	return err
}

func isTestNamespace(namespace string) bool {
	if len(namespace) < len(testNamespacePrefix) {
		return false
	}
	return namespace[0:len(testNamespacePrefix)] == testNamespacePrefix
}
