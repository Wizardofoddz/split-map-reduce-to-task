package splitmapreducejob

import "encoding/json"

// Destination is represents a container
// that can hold one or more task results
type Destination struct {
	Input json.RawMessage `json:"input"`
	UUID  string          `json:"uuid"`
}
