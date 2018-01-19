package dataset

// Leaf is a leaf node in a dataset tree
type Leaf struct {
	Content *Location `json:"content"`
}

// NewLeaf constructs a new leaf instance
func NewLeaf(contentAddr string) *Leaf {
	return &Leaf{
		Content: &Location{
			Address: contentAddr,
		},
	}
}
