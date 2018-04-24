package main

import (
	"encoding/json"
	"github.com/computes/go-sdk/pkg/types"
)

// JobInputTask is the task portion
// of a split/map/reduce job
type JobInputTask struct {
	Runner *types.Runner `json:"runner"`
}

// JobInput defines the input structure
type JobInput struct {
	Input  json.RawMessage `json:"input"`
	UUID   string          `json:"uuid"`
	Split  *JobInputTask   `json:"split"`
	Map    *JobInputTask   `json:"map"`
	Reduce *JobInputTask   `json:"reduce"`
}
