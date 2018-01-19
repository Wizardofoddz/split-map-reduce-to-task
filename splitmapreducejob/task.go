package splitmapreducejob

// JobTask is the task portion
// of a split/map/reduce job
type JobTask struct {
	Runner *Runner `json:"runner"`
}

// Task is the computes version of a task
type Task struct {
	TaskDefinition *Location  `json:"taskDefinition"`
	Input          *TaskInput `json:"input"`
	Status         *Location  `json:"status"`
}

// TaskInput is the input to a task
type TaskInput struct {
	Dataset *Location `json:"dataset"`
	Path    string    `json:"path"`
}
