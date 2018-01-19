package dataset

import "encoding/json"

// Change represents a change record in
// a dataset instance
type Change struct {
	Action string          `json:"action"`
	Path   string          `json:"path"`
	Value  json.RawMessage `json:"value"`
}
