package v2

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/workflow/v2/workbooks"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

// GetEchoWorkbookDefinition returns a simple workbook definition that does nothing except a simple "echo" command.
func GetEchoWorkbookDefinition(workbookName string) string {
	return fmt.Sprintf(`---
version: '2.0'

name: %s
description: My workbook

workflows:
  test:
    description: Simple workflow example
    type: direct
    input:
      - msg

    tasks:
      test:
        action: std.echo output="<%% $.msg %%>"`, workbookName)
}

// CreateWorkbook creates a workbook on Mistral API.
// The created workbook is a dummy workbook that performs a simple echo.
func CreateWorkbook(t *testing.T, client *gophercloud.ServiceClient) (*workbooks.Workbook, error) {
	workbookName := tools.RandomString("workbook_echo_", 5)

	definition := GetEchoWorkbookDefinition(workbookName)

	t.Logf("Attempting to create workbook: %s", workbookName)

	opts := &workbooks.CreateOpts{
		Namespace:  "some-namespace",
		Scope:      "private",
		Definition: strings.NewReader(definition),
	}

	workbook, err := workbooks.Create(context.TODO(), client, opts).Extract()
	if err != nil {
		return nil, err
	}

	t.Logf("Workbook created: %s", workbookName)

	th.AssertEquals(t, workbookName, workbook.Name)

	return workbook, nil
}

// DeleteWorkbook deletes the given workbook.
func DeleteWorkbook(t *testing.T, client *gophercloud.ServiceClient, workbook *workbooks.Workbook) {
	err := workbooks.Delete(context.TODO(), client, workbook.ID).ExtractErr()
	if err != nil {
		t.Fatalf("Unable to delete workbooks %s: %v", workbook.Name, err)
	}

	t.Logf("Deleted workbook: %s", workbook.Name)
}

// GetWorkbook gets a workbook.
func GetWorkbook(t *testing.T, client *gophercloud.ServiceClient, id string) (*workbooks.Workbook, error) {
	workbook, err := workbooks.Get(context.TODO(), client, id).Extract()
	if err != nil {
		t.Fatalf("Unable to get workbook %s: %v", id, err)
	}
	t.Logf("Workbook get: %s", workbook.Name)
	return workbook, err
}

// ListWorkbooks lists the workbooks.
func ListWorkbooks(t *testing.T, client *gophercloud.ServiceClient, opts workbooks.ListOptsBuilder) ([]workbooks.Workbook, error) {
	allPages, err := workbooks.List(client, opts).AllPages(context.TODO())
	if err != nil {
		t.Fatalf("Unable to list workbooks: %v", err)
	}
	workbooksList, err := workbooks.ExtractWorkbooks(allPages)
	if err != nil {
		t.Fatalf("Unable to extract workbooks: %v", err)
	}
	t.Logf("Workbooks list find, length: %d", len(workbooksList))
	return workbooksList, err
}
