package workbooks

import (
	"context"
	"io"

	"github.com/gophercloud/gophercloud/v2"
)

// CreateOptsBuilder allows extension to add additional parameters to the Create request.
type CreateOptsBuilder interface {
	ToWorkbookCreateParams() (io.Reader, string, error)
}

// CreateOpts specifies parameters used to create a cron trigger.
type CreateOpts struct {
	// Scope is the scope of the workbook.
	// Allowed values are "private" and "public".
	Scope string `q:"scope"`

	// Namespace will define the namespace of the workbook.
	Namespace string `q:"namespace"`

	// Definition is the workbook definition written in Mistral Workflow Language v2.
	Definition io.Reader
}

// ToWorkbookCreateParams constructs a request query string from CreateOpts.
func (opts CreateOpts) ToWorkbookCreateParams() (io.Reader, string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return opts.Definition, q.String(), err
}

// Create requests the creation of a new execution.
func Create(ctx context.Context, client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	url := createURL(client)
	var b io.Reader
	if opts != nil {
		tmpB, query, err := opts.ToWorkbookCreateParams()
		if err != nil {
			r.Err = err
			return
		}
		url += query
		b = tmpB
	}

	resp, err := client.Post(ctx, url, b, &r.Body, &gophercloud.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "text/plain",
			"Accept":       "", // Drop default JSON Accept header
		},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
