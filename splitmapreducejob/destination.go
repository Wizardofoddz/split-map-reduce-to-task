package splitmapreducejob

// Destination is represents a container
// that can hold one or more task results
type Destination struct {
	Input *Location `json:"input"`
	UUID  string    `json:"uuid"`
}
