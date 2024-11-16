/*
package workbooks provides interaction with the workbooks API in the OpenStack Mistral service.

Workbook represents a process that can be described in a various number of ways and that can do some job interesting to the end user.
Each workbook consists of tasks (at least one) describing what exact steps should be made during workbook execution.

Workbook definition is written in Mistral Workflow Language v2. You can find all specification here: https://docs.openstack.org/mistral/latest/user/wf_lang_v2.html

List workbooks

	listOpts := workbooks.ListOpts{
		Namespace: "some-namespace",
	}

	allPages, err := workbooks.List(mistralClient, listOpts).AllPages(context.TODO())
	if err != nil {
		panic(err)
	}

	allWorkbooks, err := workbooks.ExtractWorkbooks(allPages)
	if err != nil {
		panic(err)
	}

	for _, workbook := range allWorkbooks {
		fmt.Printf("%+v\n", workbook)
	}

Get a workbook

	workbook, err := workbooks.Get(context.TODO(), mistralClient, "604a3a1e-94e3-4066-a34a-aa56873ef236").Extract()
	if err != nil {
		t.Fatalf("Unable to get workbook %s: %v", id, err)
	}

	fmt.Printf("%+v\n", workbook)

Create a workbook

		workbookDefinition := `---
	      version: '2.0'

	      workbook_echo:
	        description: Simple workbook example
	        type: direct
	        input:
	          - msg

	        tasks:
	          test:
	            action: std.echo output="<% $.msg %>"`

		createOpts := &workbooks.CreateOpts{
			Definition: strings.NewReader(workbookDefinition),
			Scope: "private",
			Namespace: "some-namespace",
		}

		workbook, err := workbooks.Create(context.TODO(), mistralClient, createOpts).Extract()
		if err != nil {
			panic(err)
		}

		fmt.Printf("%+v\n", workbook)

Delete a workbook

	res := workbooks.Delete(fake.ServiceClient(), "604a3a1e-94e3-4066-a34a-aa56873ef236")
	if res.Err != nil {
		panic(res.Err)
	}
*/
package workbooks
