package client

import (
	"os"
	"strings"

	aznetwork "github.com/Azure/azure-sdk-for-go/services/network/mgmt/2017-09-01/network"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// AzureTestClient configs Azure specific clients
type AzureTestClient struct {
	networkClient aznetwork.BaseClient
}

// CreateAzureTestClient makes a new AzureTestClient
// Only consider PublicCloud Environment
func CreateAzureTestClient() (*AzureTestClient, error) {
	authconfig := AzureAuthConfigFromTestProfile()
	servicePrincipleToken, err := GetServicePrincipalToken(authconfig, &azure.PublicCloud)
	if err != nil {
		return nil, err
	}
	baseClient := aznetwork.NewWithBaseURI(azure.PublicCloud.TokenAudience, authconfig.SubscriptionID)
	baseClient.Authorizer = autorest.NewBearerAuthorizer(servicePrincipleToken)

	c := &AzureTestClient{
		networkClient: baseClient,
	}

	return c, nil
}

// CreateSubnetsClient generates subnet client with the same baseclient as azure test client
func (tc *AzureTestClient) CreateSubnetsClient() *aznetwork.SubnetsClient {
	return &aznetwork.SubnetsClient{BaseClient: tc.networkClient}
}

// CreateVirtualNetworksClient generates virtual network client with the same baseclient as azure test client
func (tc *AzureTestClient) CreateVirtualNetworksClient() *aznetwork.VirtualNetworksClient {
	return &aznetwork.VirtualNetworksClient{BaseClient: tc.networkClient}
}

func parseEnvFromLocation() *azure.Environment {
	location := os.Getenv(clusterLocationEnv)
	if strings.Contains(location, "ch") {
		return &azure.ChinaCloud
	} else if strings.Contains(location, "ger") {
		return &azure.GermanCloud
	} else if strings.Contains(location, "gov") {
		return &azure.USGovernmentCloud
	}
	return &azure.PublicCloud
}
