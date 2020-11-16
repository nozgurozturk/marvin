package parsers

import (
	"strings"
)

type Npm struct{}

func (n *Npm) Parse(file map[string]interface{}) map[string]string {

	packages := make(map[string]string)

	if file["devDependencies"] != nil {
		for key, value := range file["devDependencies"].(map[string]interface{}) {
			packages[key] = strings.Replace(value.(string), "^", "", -1)
		}
	}

	if  file["dependencies"] != nil {
		for key, value := range file["dependencies"].(map[string]interface{}) {
			packages[key] = strings.Replace(value.(string), "^", "", -1)
		}
	}

	return packages
}
