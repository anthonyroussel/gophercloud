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
description: Simple workbook example

workflows:
  workflow_echo:
    description: Simple workflow example
    type: direct
    tags:
      - tag1
      - tag2

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
