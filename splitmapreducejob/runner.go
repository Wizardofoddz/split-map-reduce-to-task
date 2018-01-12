package splitmapreducejob

import "encoding/json"

// Runner defines how a task
// should be run
type Runner struct {
	Metadata json.RawMessage `json:"metadata"`
	Type     string          `json:"type"`
}
