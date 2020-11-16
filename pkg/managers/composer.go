package managers

import (
	"encoding/json"
	"fmt"
	"github.com/nozgurozturk/marvin/client"
	"strings"
)

type Composer struct {
	apiUrl string
}

func (p *Composer) GetRegistryVersion(registryName string) (string, error) {

	// If registry is not contain owner pass
	// For exp. "php": 7.0
	if !strings.Contains(registryName, "/") {
		return "", nil
	}

	endpoint := fmt.Sprintf("/packages/%s.json", registryName)

	registryData, err := client.New(p.apiUrl).Get(endpoint, nil)
	if err != nil {
		return "", err
	}
	var registry map[string]interface{}
	if err := json.Unmarshal(registryData, &registry); err != nil {
		return "", err
	}

	var cMajor, cMinor, cPatch int
	var major, minor, patch int

	pkg := registry["package"].(map[string]interface{})
	registryVersions := pkg["versions"].(map[string]interface{})

	for _, ver := range registryVersions {

		normalizedVersion := ver.(map[string]interface{})["version_normalized"]

		// If normalized version of registry contains - or dev pass
		if strings.Contains(normalizedVersion.(string), "-") || strings.Contains(normalizedVersion.(string), "dev") {
			continue
		}

		_, err := fmt.Sscanf(normalizedVersion.(string), "%d.%d.%d", &cMajor, &cMinor, &cPatch)
		if err != nil {
			continue
		}

		// Find latest version of registry
		if cMajor > major {
			major = cMajor
			minor = cMinor
			patch = cPatch
		} else if cMinor > minor {
			minor = cMinor
			patch = cPatch
		} else if cPatch > patch {
			patch = cPatch
		}
	}

	// Converts version numbers to string
	registryVersion := fmt.Sprintf("%d.%d.%d", major, minor, patch)
	return registryVersion, nil

}
