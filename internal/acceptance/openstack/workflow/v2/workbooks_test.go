//go:build acceptance || workflow || workbooks

package v2

import (
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/workflow/v2/workbooks"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestWorkbooksCreateGetDelete(t *testing.T) {
	client, err := clients.NewWorkflowV2Client()
	th.AssertNoErr(t, err)

	workbook, err := CreateWorkbook(t, client)
	th.AssertNoErr(t, err)
	defer DeleteWorkbook(t, client, workbook)

	workbookget, err := GetWorkbook(t, client, workbook.ID)
	th.AssertNoErr(t, err)

	tools.PrintResource(t, workbookget)
}

func TestWorkbooksList(t *testing.T) {
	client, err := clients.NewWorkflowV2Client()
	th.AssertNoErr(t, err)
	workbook, err := CreateWorkbook(t, client)
	th.AssertNoErr(t, err)
	defer DeleteWorkbook(t, client, workbook)
	list, err := ListWorkbooks(t, client, &workbooks.ListOpts{
		Name: &workbooks.ListFilter{
			Value: workbook.Name,
		},
		Tags: []string{"tag1"},
		CreatedAt: &workbooks.ListDateFilter{
			Filter: workbooks.FilterGT,
			Value:  time.Now().AddDate(-1, 0, 0),
		},
	})
	th.AssertNoErr(t, err)
	th.AssertEquals(t, 1, len(list))
	tools.PrintResource(t, list)
}
