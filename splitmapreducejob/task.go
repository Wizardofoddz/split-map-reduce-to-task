package splitmapreducejob

// JobTask is the task portion
// of a split/map/reduce job
type JobTask struct {
	Runner *Runner `json:"runner"`
}
