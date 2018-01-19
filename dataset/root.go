package dataset

// Root is the root node of a dataset structure
type Root struct {
	Left *Location `json:"left"`
}

// NewRoot constructs a new Root instance
func NewRoot(leftAddr string) *Root {
	return &Root{
		Left: &Location{
			Address: leftAddr,
		},
	}
}
