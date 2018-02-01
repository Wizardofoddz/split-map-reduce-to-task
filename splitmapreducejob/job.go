package splitmapreducejob

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/computes/ipfs-http-api/dag"
	"github.com/computes/split-map-reduce-to-task/dataset"
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
	return dataset.Create(job.ipfsURL, &dataset.Change{
		Action: "set",
		Path:   "split/input",
		Value:  job.Input,
	})
}

// StoreSplitTask stores the split task
func (job *Job) StoreSplitTask() (string, error) {
	destinationAddr, err := job.StoreDestination()
	if err != nil {
		return "", err
	}

	taskStatusAddr, err := job.StoreTaskStatus()
	if err != nil {
		return "", err
	}

	defAddrs, err := job.StoreTaskDefinitions(destinationAddr)
	if err != nil {
		return "", err
	}

	splitDefAddr := defAddrs[0]

	taskBytes, err := json.Marshal(&Task{
		Input: &TaskInput{
			Dataset: &Location{Address: destinationAddr},
			Path:    "split/input",
		},
		Status: &Location{
			Address: taskStatusAddr,
		},
		TaskDefinition: &Location{Address: splitDefAddr},
	})
	if err != nil {
		return "", err
	}

	return dag.Put(job.ipfsURL, bytes.NewReader(taskBytes))
}

// StoreTaskDefinitions will store the 3 task definitions
// and return their addresses
func (job *Job) StoreTaskDefinitions(destinationAddr string) ([]string, error) {
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

// StoreTaskStatus creates a task status dataset
// and writes the job UUID to it
func (job *Job) StoreTaskStatus() (string, error) {
	uuidBytes, err := json.Marshal(job.UUID)
	if err != nil {
		return "", err
	}

	uuidAddr, err := dag.Put(job.ipfsURL, bytes.NewBuffer(uuidBytes))
	if err != nil {
		return "", err
	}

	valueBytes, err := json.Marshal(&Location{Address: uuidAddr})
	if err != nil {
		return "", err
	}

	return dataset.Create(job.ipfsURL, &dataset.Change{
		Action: "set",
		Path:   "uuid",
		Value:  json.RawMessage(valueBytes),
	})
}

func (job *Job) mapTaskDefinition(destinationAddr, reduceAddr string) *TaskDefinition {
	return &TaskDefinition{
		Runner: job.Map.Runner,
		Result: &TaskDefinitionResult{
			Action: "append",
			Destination: &DataRef{
				Dataset: &Location{
					Address: destinationAddr,
				},
				Path: "map/results",
			},
		},
		Conditions: []*TaskDefinitionCondition{
			&TaskDefinitionCondition{
				Name: "Create a Reduce Task",
				Condition: fmt.Sprintf(
					"len(dataset(hpcp('%v/split/results'))) == len(dataset(hpcp('%v/map/results')))  && !exist(dataset(hpcp('%v/reduce/results')))",
					destinationAddr,
					destinationAddr,
					destinationAddr,
				),
				TaskDefinition: &Location{Address: reduceAddr},
				Source: &DataRef{
					Dataset: &Location{Address: destinationAddr},
					Path:    "map/results",
				},
			},
		},
	}
}

func (job *Job) reduceTaskDefinition(destinationAddr string) *TaskDefinition {
	return &TaskDefinition{
		Runner: job.Reduce.Runner,
		Result: &TaskDefinitionResult{
			Action: "set",
			Destination: &DataRef{
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
			Destination: &DataRef{
				Dataset: &Location{
					Address: destinationAddr,
				},
				Path: "split/results",
			},
		},
		Conditions: []*TaskDefinitionCondition{
			&TaskDefinitionCondition{
				Name: "Create Map Tasks",
				Condition: fmt.Sprintf(
					"exist(dataset(hpcp('%v/split/results'))) && len(dataset(hpcp('%v/map/results'))) < len(dataset(hpcp('%v/split/results')))",
					destinationAddr,
					destinationAddr,
					destinationAddr,
				),
				TaskDefinition: &Location{
					Address: mapAddr,
				},
				Action: "map",
				Source: &DataRef{
					Dataset: &Location{Address: destinationAddr},
					Path:    "split/results",
				},
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
