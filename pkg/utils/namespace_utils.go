/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package utils

import (
	"fmt"
	"strings"
	"time"

	"k8s.io/api/core/v1"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilnet "k8s.io/apimachinery/pkg/util/net"
	"k8s.io/apimachinery/pkg/util/wait"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

const (
	deletionTimeout   = 10 * time.Minute
	poll              = 2 * time.Second
	singleCallTimeout = 5 * time.Minute
)

func findExistingKubeConfig() string {
	// locations using DefaultClientConfigLoadingRules
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	return rules.GetDefaultFilename()
}

//CreateTestingNamespace builds namespace for each test
//baseName and labels determine name of the space
func CreateTestingNamespace(baseName string, cs clientset.Interface) (*v1.Namespace, error) {
	Logf("Creating a test namespace")

	namespaceObj := &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: fmt.Sprintf("e2e-tests-%v-", baseName),
			Namespace:    "",
		},
		Status: v1.NamespaceStatus{},
	}
	// Be robust about making the namespace creation call.
	var got *v1.Namespace
	if err := wait.PollImmediate(poll, 30*time.Second, func() (bool, error) {
		var err error
		got, err = cs.CoreV1().Namespaces().Create(namespaceObj)
		if err != nil {
			if IsRetryableAPIError(err) {
				return false, nil
			}
			return false, err
		}
		return true, nil
	}); err != nil {
		return nil, err
	}
	return got, nil
}

// GetNamespacesList gets the namespace list
func GetNamespaceList(cs clientset.Interface) (*v1.NamespaceList, error) {
	var list *v1.NamespaceList
	if err := wait.PollImmediate(poll, 30*time.Second, func() (bool, error) {
		var err error
		list, err = cs.CoreV1().Namespaces().List(metav1.ListOptions{})
		if err != nil {
			if IsRetryableAPIError(err) {
				return false, nil
			}
			return false, err
		}
		return true, nil
	}); err != nil {
		return nil, err
	}
	return list, nil
}

// DeleteNamespace deletes the provided namespace, waits for it to be completely deleted, and then checks
// whether there are any pods remaining in a non-terminating state.
func DeleteNamespace(cs clientset.Interface, namespace string) error {
	Logf("Deleting namespace %s", namespace)
	if err := cs.CoreV1().Namespaces().Delete(namespace, nil); err != nil {
		return err
	}
	// wait for namespace to delete or timeout.
	err := wait.PollImmediate(poll, deletionTimeout, func() (bool, error) {
		if _, err := cs.CoreV1().Namespaces().Get(namespace, metav1.GetOptions{}); err != nil {
			if apierrs.IsNotFound(err) {
				return true, nil
			}
			Logf("Error while waiting for namespace to be terminated: %v", err)
			if !IsRetryableAPIError(err) {
				return false, err
			}
		}
		return false, nil
	})
	return err
}

// IsRetryableAPIError will judge whether an error retrable or not
func IsRetryableAPIError(err error) bool {
	// These errors may indicate a transient error that we can retry in tests.
	if apierrs.IsInternalError(err) || apierrs.IsTimeout(err) || apierrs.IsServerTimeout(err) ||
		apierrs.IsTooManyRequests(err) || utilnet.IsProbableEOF(err) || utilnet.IsConnectionReset(err) {
		return true
	}
	// If the error sends the Retry-After header, we respect it as an explicit confirmation we should retry.
	if _, shouldRetry := apierrs.SuggestsClientDelay(err); shouldRetry {
		return true
	}
	return false
}

// GetResourceGroup get RG name which is same of cluster name as definited in k8s-azure
func GetResourceGroup() string {
	return ExtractDNSPrefix()
}

// ExtractDNSPrefix obtains the cluster DNS prefix
func ExtractDNSPrefix() string {
	c := obtainConfig()
	return c.CurrentContext
}

// extractSuffix obtains the server domain name suffix
func extractSuffix() string {
	c := obtainConfig()
	prefix := ExtractDNSPrefix()
	url := c.Clusters[prefix].Server
	suffix := url[strings.Index(url, "."):]
	return suffix
}

// Load config from file
func obtainConfig() *clientcmdapi.Config {
	filename := findExistingKubeConfig()
	c := clientcmd.GetConfigFromFileOrDie(filename)
	return c
}

// stringInSlice check if string in a list
func stringInSlice(s string, list []string) bool {
	for _, item := range list {
		if item == s {
			return true
		}
	}
	return false
}
