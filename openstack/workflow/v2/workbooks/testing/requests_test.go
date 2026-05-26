package testing

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2/openstack/workflow/v2/workbooks"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	fake "github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestCreateWorkbook(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	definition := `---
version: '2.0'

name: workbook_echo
description: Simple workbook example
tags:
  - workbook_tag1
  - workbook_tag2

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
        action: std.echo output="<% $.msg %>"
`

	th.Mux.HandleFunc("/workbooks", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "text/plain")
		th.TestFormValues(t, r, map[string]string{
			"namespace": "some-namespace",
			"scope":     "private",
		})
		th.TestBody(t, r, definition)

		w.WriteHeader(http.StatusCreated)
		w.Header().Add("Content-Type", "application/json")

		fmt.Fprintf(w, `{
			"created_at": "2024-11-16 22:48:45",
			"definition": "---\nversion: '2.0'\n\nname: workbook_echo\ndescription: Simple workbook example\ntags:\n  - workbook_tag1\n  - workbook_tag2\n\nworkflows:\n  workflow_echo:\n    description: Simple workflow example\n    type: direct\n    tags:\n      - tag1\n      - tag2\n\n    input:\n      - msg\n\n    tasks:\n      test:\n        action: std.echo output=\"<%% $.msg %%>\"",
			"id": "c20698c4-2ba4-4334-9c52-7f0d3b0af21f",
			"name": "workbook_echo",
			"namespace": "some-namespace",
			"project_id": "778c0f25df0d492a9a868ee9e2fbb513",
			"scope": "private",
			"tags": ["workbook_tag1", "workbook_tag2"],
			"updated_at": "2024-11-16 22:48:45"
		}`)
	})

	opts := &workbooks.CreateOpts{
		Namespace:  "some-namespace",
		Scope:      "private",
		Definition: strings.NewReader(definition),
	}

	actual, err := workbooks.Create(context.TODO(), fake.ServiceClient(), opts).Extract()
	if err != nil {
		t.Fatalf("Unable to create workbook: %v", err)
	}

	updated := time.Date(2024, time.November, 16, 22, 48, 45, 0, time.UTC)
	expected := &workbooks.Workbook{
		ID:         "c20698c4-2ba4-4334-9c52-7f0d3b0af21f",
		Definition: "---\nversion: '2.0'\n\nname: workbook_echo\ndescription: Simple workbook example\ntags:\n  - workbook_tag1\n  - workbook_tag2\n\nworkflows:\n  workflow_echo:\n    description: Simple workflow example\n    type: direct\n    tags:\n      - tag1\n      - tag2\n\n    input:\n      - msg\n\n    tasks:\n      test:\n        action: std.echo output=\"<% $.msg %>\"",
		Name:       "workbook_echo",
		Namespace:  "some-namespace",
		ProjectID:  "778c0f25df0d492a9a868ee9e2fbb513",
		Scope:      "private",
		Tags:       []string{"workbook_tag1", "workbook_tag2"},
		CreatedAt:  time.Date(2024, time.November, 16, 22, 48, 45, 0, time.UTC),
		UpdatedAt:  &updated,
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %#v, but was %#v", expected, actual)
	}
}
