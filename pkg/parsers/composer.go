package parsers

import (
	"strings"
)

type Composer struct{}

func (c *Composer) Parse(file map[string]interface{}) map[string]string {

	packages := make(map[string]string)

	for key, value := range file["require-dev"].(map[string]interface{}) {
		packages[key] = getLatestVersion(value.(string))
	}

	for key, value := range file["require"].(map[string]interface{}) {
		packages[key] = getLatestVersion(value.(string))
	}

	return packages
}

// TODO: improve this detection

// Gets latest version of package version
func getLatestVersion(rawVersion string) string {
	versions := strings.Split(
		strings.Replace(
			strings.Replace(rawVersion, "*", "", -1), "^", "", -1),
		"|")
	latest := versions[len(versions)-1]
	return latest
}
