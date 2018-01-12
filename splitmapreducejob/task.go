package splitmapreducejob

// JobTask is the task portion
// of a split/map/reduce job
type JobTask struct {
	Runner *Runner `json:"runner"`
}

// Task is the computes version of a task
type Task struct {
	TaskDefinition *Location `json:"taskDefinitiion"`
	Input          *Location `json:"input"`
}
