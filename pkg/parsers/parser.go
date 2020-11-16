/*
	Package parsers gets packages with name and version in package files and parses
 */
package parsers

import (
	"errors"
	"fmt"
)

const (
	npm      = "package.json"
	composer = "composer.json"
)

type Parser interface {
	Parse(file map[string]interface{}) map[string]string // Parses package file
}

// Create Parser with given package file name
func NewParser(packageFileName string) (Parser, error) {
	switch packageFileName {
	case npm:
		return new(Npm), nil
	case composer:
		return new(Composer), nil
	default:
		return nil, errors.New(fmt.Sprintf("Undefined package file: %s", packageFileName))
	}
}
