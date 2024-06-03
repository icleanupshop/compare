// tfvars provides the ability to read unstructured data from a tfvars file
package tfvars

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/hashicorp/hcl"
)

// Tfvars holds the data from a tfvars file
type Tfvars struct {
	path string
	data map[string]interface{}
}

// New creates a new Tfvars object
func New(path string) (*Tfvars, error) {
	var err error
	t := &Tfvars{path: path, data: nil}
	t.data, err = t.read()

	return t, err
}

// Get returns the string value of the object at the given path
func (t *Tfvars) Get(arg ...string) string {
	data := t.Raw(arg...)

	if v, ok := data.(string); ok {
		return v
	}

	return fmt.Sprintf("%v", data)
}

// Keys returns the keys available at the given path
func (t *Tfvars) Keys(arg ...string) []string {
	data := t.Raw(arg...)

	return t.keys(data)
}

// Raw returns the raw interface{} value at the given path
func (t *Tfvars) Raw(arg ...string) interface{} {
	var data interface{}
	data = t.data

	for _, a := range arg {
		data = t.get(data, a)
		if data == nil {
			return data
		}
	}

	return data
}

func (t *Tfvars) get(data interface{}, arg string) interface{} {
	switch tdata := data.(type) {
	case map[string]interface{}:
		if tdata[arg] != nil {
			return tdata[arg]
		}
	case []map[string]interface{}:
		for _, d := range tdata {
			if v, ok := d[arg]; ok {
				return v
			}
		}
	}

	return nil
}

func (t *Tfvars) keys(data interface{}) []string {
	var keys []string
	switch tdata := data.(type) {
	case map[string]interface{}:
		for k := range tdata {
			keys = append(keys, k)
		}
	case []map[string]interface{}:
		for _, d := range tdata {
			for k := range d {
				keys = append(keys, k)
			}
		}
	}

	return keys
}

func (t *Tfvars) read() (map[string]interface{}, error) {
	var result map[string]interface{}
	f, err := os.Open(t.path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %s -> %w", t.path, err)
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %s -> %w", t.path, err)
	}
	err = hcl.Decode(&result, string(b))
	if err != nil {
		return result, fmt.Errorf("failed to decode hcl file: %s -> %w", t.path, err)
	}

	return result, nil
}
