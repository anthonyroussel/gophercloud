package workbooks

import (
	"encoding/json"
	"time"

	"github.com/gophercloud/gophercloud/v2"
)

// CreateResult is the response of a Post operations. Call its Extract method to interpret it as a list of Workbooks.
type CreateResult struct {
	gophercloud.Result
}

// Extract helps to get created Workbook struct from a Create function.
func (r CreateResult) Extract() (*Workbook, error) {
	var s Workbook
	err := r.ExtractInto(&s)
	return &s, err
}

// Workbook represents a workbook execution on OpenStack mistral API.
type Workbook struct {
	// ID is the workbook's unique ID.
	ID string `json:"id"`

	// Definition is the workbook definition in Mistral v2 DSL.
	Definition string `json:"definition"`

	// Name is the name of the workbook.
	Name string `json:"name"`

	// Namespace is the namespace of the workbook.
	Namespace string `json:"namespace"`

	// ProjectID is the project id owner of the workbook.
	ProjectID string `json:"project_id"`

	// Scope is the scope of the workbook.
	// Values can be "private" or "public".
	Scope string `json:"scope"`

	// Tags is a list of tags associated to the workbook.
	Tags []string `json:"tags"`

	// CreatedAt is the creation date of the workbook.
	CreatedAt time.Time `json:"-"`

	// UpdatedAt is the last update date of the workbook.
	UpdatedAt *time.Time `json:"-"`
}

// UnmarshalJSON implements unmarshalling custom types
func (r *Workbook) UnmarshalJSON(b []byte) error {
	type tmp Workbook
	var s struct {
		tmp
		CreatedAt gophercloud.JSONRFC3339ZNoTNoZ  `json:"created_at"`
		UpdatedAt *gophercloud.JSONRFC3339ZNoTNoZ `json:"updated_at"`
	}

	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	*r = Workbook(s.tmp)

	r.CreatedAt = time.Time(s.CreatedAt)
	if s.UpdatedAt != nil {
		t := time.Time(*s.UpdatedAt)
		r.UpdatedAt = &t
	}

	return nil
}
