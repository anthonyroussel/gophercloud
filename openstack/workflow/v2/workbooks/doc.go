/*
package workbooks provides interaction with the workbooks API in the OpenStack Mistral service.

Workbook represents a process that can be described in a various number of ways and that can do some job interesting to the end user.
Each workbook consists of tasks (at least one) describing what exact steps should be made during workbook execution.

Workbook definition is written in Mistral Workflow Language v2. You can find all specification here: https://docs.openstack.org/mistral/latest/user/wf_lang_v2.html

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
*/
package workbooks
