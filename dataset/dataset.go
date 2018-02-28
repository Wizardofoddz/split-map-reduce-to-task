package dataset

import (
	"bytes"
	"encoding/json"
	"net/url"

	"github.com/computes/ipfs-http-api/dag"
)

// Create constructs and stores an dataset given
// an initial change.
func Create(ipfsURL url.URL, initialChange *Change) (string, error) {
	initialChangeAddr, err := dagPut(ipfsURL, initialChange)
	if err != nil {
		return "", err
	}

	node := _Node{
		Contents: []_Ref{
			_Ref{Address: initialChangeAddr},
		},
	}
	return dagPut(ipfsURL, node)
}

func dagPut(ipfsURL url.URL, obj interface{}) (string, error) {
	data, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}

	return dag.Put(ipfsURL, bytes.NewBuffer(data))
}

type _Node struct {
	Contents []_Ref `json:"contents"`
}

type _Ref struct {
	Address string `json:"/"`
}
