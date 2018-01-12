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

// New constructs a new Job instance
func New(ipfsURL url.URL) *Job {
	return &Job{ipfsURL: ipfsURL}
}

// StoreDestination stores the destination and returns
// the location
func (job *Job) StoreDestination() (string, error) {
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

// StoreTaskDefinitions will store the 3 task definitions
// and return their addresses
func (job *Job) StoreTaskDefinitions() ([]string, error) {
	destinationAddr, err := job.StoreDestination()
	if err != nil {
		return nil, err
	}

	reduceAddr, err := job.storeReduceTaskDefinition(destinationAddr)
	if err != nil {
		return nil, err
	}

	mapAddr, err := job.storeMapTaskDefinition(destinationAddr, reduceAddr)
	if err != nil {
		return nil, err
	}

	splitAddr, err := job.storeSplitTaskDefinition(destinationAddr, mapAddr)
	if err != nil {
		return nil, err
	}

	return []string{splitAddr, mapAddr, reduceAddr}, nil
}

func (job *Job) mapTaskDefinition(destinationAddr, reduceAddr string) *TaskDefinition {
	return &TaskDefinition{
		Runner: job.Map.Runner,
		Result: &TaskDefinitionResult{
			Action: "append",
			Destination: &TaskDefinitionResultDestination{
				Dataset: &Location{
					Address: destinationAddr,
				},
				Path: "map/results",
			},
		},
		Conditions: []*TaskDefinitionCondition{
			&TaskDefinitionCondition{
				Name: "Make a Reduce Task",
				Rule: "len(split/results) == len(map/results)",
				TaskDefinition: &Location{
					Address: reduceAddr,
				},
				Map:    false,
				Source: "map/results",
			},
		},
	}
}

func (job *Job) reduceTaskDefinition(destinationAddr string) *TaskDefinition {
	return &TaskDefinition{
		Runner: job.Reduce.Runner,
		Result: &TaskDefinitionResult{
			Action: "set",
			Destination: &TaskDefinitionResultDestination{
				Dataset: &Location{
					Address: destinationAddr,
				},
				Path: "reduce/results",
			},
		},
		Conditions: nil,
	}
}

func (job *Job) splitTaskDefinition(destinationAddr, mapAddr string) *TaskDefinition {
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
		Conditions: []*TaskDefinitionCondition{
			&TaskDefinitionCondition{
				Name: "Make Map Tasks",
				Rule: "len(split/results) > 0",
				TaskDefinition: &Location{
					Address: mapAddr,
				},
				Map:    true,
				Source: "split/results",
			},
		},
	}
}

func (job *Job) storeMapTaskDefinition(destinationAddr, reduceAddr string) (string, error) {
	definition := job.mapTaskDefinition(destinationAddr, reduceAddr)

	data, err := json.Marshal(definition)
	if err != nil {
		return "", err
	}

	return dag.Put(job.ipfsURL, bytes.NewBuffer(data))
}

func (job *Job) storeReduceTaskDefinition(destinationAddr string) (string, error) {
	definition := job.reduceTaskDefinition(destinationAddr)

	data, err := json.Marshal(definition)
	if err != nil {
		return "", err
	}

	return dag.Put(job.ipfsURL, bytes.NewBuffer(data))
}

func (job *Job) storeSplitTaskDefinition(destinationAddr, mapAddr string) (string, error) {
	definition := job.splitTaskDefinition(destinationAddr, mapAddr)

	data, err := json.Marshal(definition)
	if err != nil {
		return "", err
	}

	return dag.Put(job.ipfsURL, bytes.NewBuffer(data))
}
