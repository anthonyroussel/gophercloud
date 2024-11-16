package workbooks

import (
	"github.com/gophercloud/gophercloud/v2"
)

func createURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("workbooks")
}

func deleteURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL("workbooks", id)
}

func getURL(client *gophercloud.ServiceClient, id string) string {
	return client.ServiceURL("workbooks", id)
}

func listURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("workbooks")
}
