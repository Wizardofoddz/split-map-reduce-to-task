package dataset

import (
	"bytes"
	"encoding/json"
	"net/url"

	TREE "github.com/computes/fake-ipfs-merkle-tree"
	"github.com/computes/ipfs-http-api/dag"
)

// Create constructs and stores an dataset given
// an initial change.
func Create(ipfsURL url.URL, initialChange *Change) (string, error) {
	initialChangeAddr, err := dagPut(ipfsURL, initialChange)
	if err != nil {
		return "", err
	}

	tree := TREE.FromContentAddresses(ipfsURL, []string{initialChangeAddr})
	return tree.Store()
}

func dagPut(ipfsURL url.URL, obj interface{}) (string, error) {
	data, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}

	return dag.Put(ipfsURL, bytes.NewBuffer(data))
}
