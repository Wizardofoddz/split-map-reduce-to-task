package splitmapreducejob

import (
	"bytes"
	"encoding/json"
	"net/url"

	"github.com/computes/ipfs-http-api/dag"
)

// Job represents a standard
// split/map/reduce job and can be converted
// to a computes task set. This struct should
// be instantiated using New so that ipfsURL may
// be set
type Job struct {
	ipfsURL url.URL
	Input   json.RawMessage `json:"input"`
	UUID    string          `json:"uuid"`
	Split   *JobTask        `json:"split"`
	Map     *JobTask        `json:"map"`
	Reduce  *JobTask        `json:"reduce"`
}

// InitializeDestination stores the destination and returns
// the location
func (job *Job) InitializeDestination() (string, error) {
	destination := &Destination{
		Input: job.Input,
		UUID:  job.UUID,
	}

	destinationBytes, err := json.Marshal(destination)
	if err != nil {
		return "", err
	}

	return dag.Put(job.ipfsURL, bytes.NewBuffer(destinationBytes))
}

// New constructs a new Job instance
func New(ipfsURL url.URL) *Job {
	return &Job{ipfsURL: ipfsURL}
}

// TasksDefinitions will return 3 task definitions,
// one each for split, map, and reduce
func (job *Job) TasksDefinitions() ([]*TaskDefinition, error) {
	destinationAddr, err := job.InitializeDestination()
	if err != nil {
		return nil, err
	}

	return []*TaskDefinition{
		job.splitTaskDefinition(destinationAddr),
		job.mapTaskDefinition(destinationAddr),
		job.reduceTaskDefinition(destinationAddr),
	}, nil
}

func (job *Job) mapTaskDefinition(destinationAddr string) *TaskDefinition {
	return nil
}

func (job *Job) reduceTaskDefinition(destinationAddr string) *TaskDefinition {
	return nil
}

func (job *Job) splitTaskDefinition(destinationAddr string) *TaskDefinition {
	return &TaskDefinition{
		Runner: job.Split.Runner,
		Result: &TaskDefinitionResult{
			Action: "set",
			Destination: &TaskDefinitionResultDestination{
				Dataset: &Location{
					Address: destinationAddr,
				},
				Path: "split/results",
			},
		},
	}
}
