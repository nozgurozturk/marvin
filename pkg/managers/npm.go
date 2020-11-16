package managers

import (
	"encoding/json"
	"fmt"
	"github.com/nozgurozturk/marvin/client"
)

type Npm struct {
	apiUrl string
}

func (n *Npm) GetRegistryVersion(registryName string) (string, error) {

	endpoint := fmt.Sprintf("/%s", registryName)
	registryData, err := client.New(n.apiUrl).Get(endpoint, nil)
	if err != nil {
		return "", err
	}
	var registry map[string]interface{}
	if err := json.Unmarshal(registryData, &registry); err != nil {
		return "", err
	}
	registryVersion := registry["dist-tags"].(map[string]interface{})["latest"].(string)
	return registryVersion, nil
}

