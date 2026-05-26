//go:build acceptance || workflow || workbooks

package v2

import (
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestWorkbooksCreateGetDelete(t *testing.T) {
	client, err := clients.NewWorkflowV2Client()
	th.AssertNoErr(t, err)

	workbook, err := CreateWorkbook(t, client)
	th.AssertNoErr(t, err)

	tools.PrintResource(t, workbook)
}
