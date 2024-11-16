package testing

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2/openstack/workflow/v2/workbooks"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	fake "github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestCreateWorkbook(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	definition := `---
version: '2.0'

name: my_workbook
description: My workbook

workflows:
  workflow_echo:
    description: Simple workflow example
    type: direct
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
			"definition": "---\nversion: '2.0'\n\nname: my_workbook\ndescription: My workbook\n\nworkflows:\n  workflow_echo:\n    description: Simple workflow example\n    type: direct\n    input:\n      - msg\n\n    tasks:\n      test:\n        action: std.echo output=\"<%% $.msg %%>\"",
			"id": "c20698c4-2ba4-4334-9c52-7f0d3b0af21f",
			"name": "my_workbook",
			"namespace": "some-namespace",
			"project_id": "778c0f25df0d492a9a868ee9e2fbb513",
			"scope": "private",
			"tags": [],
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
		Definition: "---\nversion: '2.0'\n\nname: my_workbook\ndescription: My workbook\n\nworkflows:\n  workflow_echo:\n    description: Simple workflow example\n    type: direct\n    input:\n      - msg\n\n    tasks:\n      test:\n        action: std.echo output=\"<% $.msg %>\"",
		Name:       "my_workbook",
		Namespace:  "some-namespace",
		ProjectID:  "778c0f25df0d492a9a868ee9e2fbb513",
		Scope:      "private",
		Tags:       []string{},
		CreatedAt:  time.Date(2024, time.November, 16, 22, 48, 45, 0, time.UTC),
		UpdatedAt:  &updated,
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %#v, but was %#v", expected, actual)
	}
}

func TestDeleteWorkbook(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/workbooks/c20698c4-2ba4-4334-9c52-7f0d3b0af21f", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.WriteHeader(http.StatusAccepted)
	})

	res := workbooks.Delete(context.TODO(), fake.ServiceClient(), "c20698c4-2ba4-4334-9c52-7f0d3b0af21f")
	th.AssertNoErr(t, res.Err)
}

func TestGetWorkbook(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	th.Mux.HandleFunc("/workbooks/1", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-token", fake.TokenID)
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, `
			{
				"created_at": "2024-11-16 22:48:45",
				"definition": "---\nversion: '2.0'\n\nname: my_workbook\ndescription: My workbook\n\nworkflows:\n  workflow_echo:\n    description: Simple workflow example\n    type: direct\n    input:\n      - msg\n\n    tasks:\n      test:\n        action: std.echo output=\"<%% $.msg %%>\"",
				"id": "c20698c4-2ba4-4334-9c52-7f0d3b0af21f",
				"name": "my_workbook",
				"namespace": "some-namespace",
				"project_id": "778c0f25df0d492a9a868ee9e2fbb513",
				"scope": "private",
				"tags": [],
				"updated_at": "2024-11-16 22:48:45"
			}
		`)
	})
	actual, err := workbooks.Get(context.TODO(), fake.ServiceClient(), "1").Extract()
	if err != nil {
		t.Fatalf("Unable to get workbook: %v", err)
	}

	updated := time.Date(2024, time.November, 16, 22, 48, 45, 0, time.UTC)
	expected := &workbooks.Workbook{
		ID:         "c20698c4-2ba4-4334-9c52-7f0d3b0af21f",
		Definition: "---\nversion: '2.0'\n\nname: my_workbook\ndescription: My workbook\n\nworkflows:\n  workflow_echo:\n    description: Simple workflow example\n    type: direct\n    input:\n      - msg\n\n    tasks:\n      test:\n        action: std.echo output=\"<% $.msg %>\"",
		Name:       "my_workbook",
		Namespace:  "some-namespace",
		ProjectID:  "778c0f25df0d492a9a868ee9e2fbb513",
		Scope:      "private",
		Tags:       []string{},
		CreatedAt:  time.Date(2024, time.November, 16, 22, 48, 45, 0, time.UTC),
		UpdatedAt:  &updated,
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %#v, but was %#v", expected, actual)
	}
}

func TestListWorkbooks(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	th.Mux.HandleFunc("/workbooks", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.Header().Add("Content-Type", "application/json")
		if err := r.ParseForm(); err != nil {
			t.Errorf("Failed to parse request form %v", err)
		}
		marker := r.Form.Get("marker")
		switch marker {
		case "":
			fmt.Fprintf(w, `{
				"next": "%s/workbooks?marker=c20698c4-2ba4-4334-9c52-7f0d3b0af21f",
				"workbooks": [
					{
						"created_at": "2024-11-16 22:48:45",
						"definition": "---\nversion: '2.0'\n\nname: my_workbook\ndescription: My workbook\n\nworkflows:\n  workflow_echo:\n    description: Simple workflow example\n    type: direct\n    input:\n      - msg\n\n    tasks:\n      test:\n        action: std.echo output=\"<%% $.msg %%>\"",
						"id": "c20698c4-2ba4-4334-9c52-7f0d3b0af21f",
						"name": "my_workbook",
						"namespace": "some-namespace",
						"project_id": "778c0f25df0d492a9a868ee9e2fbb513",
						"scope": "private",
						"tags": [],
						"updated_at": "2024-11-16 22:48:45"
					}
				]
			}`, th.Server.URL)
		case "c20698c4-2ba4-4334-9c52-7f0d3b0af21f":
			fmt.Fprintf(w, `{ "workbooks": [] }`)
		default:
			t.Fatalf("Unexpected marker: [%s]", marker)
		}
	})
	pages := 0
	// Get all workbooks
	err := workbooks.List(fake.ServiceClient(), nil).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		pages++
		actual, err := workbooks.ExtractWorkbooks(page)
		if err != nil {
			return false, err
		}

		updated := time.Date(2024, time.November, 16, 22, 48, 45, 0, time.UTC)
		expected := []workbooks.Workbook{
			{
				ID:         "c20698c4-2ba4-4334-9c52-7f0d3b0af21f",
				Definition: "---\nversion: '2.0'\n\nname: my_workbook\ndescription: My workbook\n\nworkflows:\n  workflow_echo:\n    description: Simple workflow example\n    type: direct\n    input:\n      - msg\n\n    tasks:\n      test:\n        action: std.echo output=\"<% $.msg %>\"",
				Name:       "my_workbook",
				Namespace:  "some-namespace",
				ProjectID:  "778c0f25df0d492a9a868ee9e2fbb513",
				Scope:      "private",
				Tags:       []string{},
				CreatedAt:  time.Date(2024, time.November, 16, 22, 48, 45, 0, time.UTC),
				UpdatedAt:  &updated,
			},
		}

		if !reflect.DeepEqual(expected, actual) {
			t.Errorf("Expected %#v, but was %#v", expected, actual)
		}
		return true, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if pages != 1 {
		t.Errorf("Expected one page, got %d", pages)
	}
}

func TestToWorkbookListQuery(t *testing.T) {
	for expected, opts := range map[string]*workbooks.ListOpts{
		newValue("tags", `tag1,tag2`): {
			Tags: []string{"tag1", "tag2"},
		},
		newValue("name", `neq:invalid_name`): {
			Name: &workbooks.ListFilter{
				Filter: workbooks.FilterNEQ,
				Value:  "invalid_name",
			},
		},
		newValue("created_at", `gt:2024-01-01 00:00:00`): {
			CreatedAt: &workbooks.ListDateFilter{
				Filter: workbooks.FilterGT,
				Value:  time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
		},
	} {
		actual, _ := opts.ToWorkbookListQuery()
		th.AssertEquals(t, expected, actual)
	}
}
func newValue(param, value string) string {
	v := url.Values{}
	v.Add(param, value)
	return "?" + v.Encode()
}
