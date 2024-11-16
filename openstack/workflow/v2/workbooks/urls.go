package workbooks

import (
	"github.com/gophercloud/gophercloud/v2"
)

func createURL(client *gophercloud.ServiceClient) string {
	return client.ServiceURL("workbooks")
}
