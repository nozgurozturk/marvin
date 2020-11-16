/*
 Package managers consumes package manager's API
*/

package managers

import (
	"errors"
	"fmt"
)

const (
	npm      = "package.json"
	composer = "composer.json"
)

type Manager interface {
	GetRegistryVersion(registryName string) (string, error) // Gets registry's version
}

// Creates new manager with given file name
func NewManager(fileName string) (Manager, error) {
	switch fileName {
	case npm:
		p := new(Npm)
		p.apiUrl = "https://registry.npmjs.org"
		return p, nil
	case composer:
		p := new(Composer)
		p.apiUrl = "https://packagist.org"
		return p, nil
	default:
		return nil, errors.New(fmt.Sprintf("Undefined package file name: %s", fileName))
	}
}
