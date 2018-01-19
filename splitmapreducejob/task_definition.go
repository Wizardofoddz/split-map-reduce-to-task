package splitmapreducejob

// TaskDefinition is the canonical computes
// task definition format
type TaskDefinition struct {
	Runner     *Runner                    `json:"runner"`
	Result     *TaskDefinitionResult      `json:"result"`
	Conditions []*TaskDefinitionCondition `json:"conditions"`
}

// TaskDefinitionCondition defines under what
// conditions new tasks are generated
type TaskDefinitionCondition struct {
	Name           string    `json:"name"`
	Condition      string    `json:"condition"`
	TaskDefinition *Location `json:"taskDefinition"`
	Action         string    `json:"action"`
	Source         *DataRef  `json:"source"`
}

// TaskDefinitionResult defines what to do with the
// result after the task is done
type TaskDefinitionResult struct {
	Action      string   `json:"action"`
	Destination *DataRef `json:"destination"`
}

// DataRef describes an IPLD ref with a path
// component
type DataRef struct {
	Dataset *Location `json:"dataset"`
	Path    string    `json:"path"`
}
