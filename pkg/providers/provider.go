/*
Package providers is resolve repository data with given providers
Accepted providers [github, gitlab]
*/
package providers

import (
	"errors"
	"fmt"
	"net/url"
)

const (
	github    = "github.com"
	gitlab    = "gitlab.com"
)

type Provider interface {
	UrlResolver() (string, string)                                                  // Gets owner and name of repository
	GetRepositoryTree(owner string, name string) ([]map[string]interface{}, error)  // Gets repository tree of main directory
	FindPackagesInfo(tree []map[string]interface{}) []map[string]interface{}        // Gets package manager file info from provider's API
	GetPackageFiles(files []map[string]interface{}) (map[string]interface{}, error) // Gets packages files with name and version
}

// Detect provider from given url
func GetProvider(u *url.URL) (Provider, error) {
	switch u.Host {
	case github:
		g := new(Github)
		g.url = u
		g.apiUrl = "https://api.github.com"
		return g, nil
	case gitlab:
		g := new(Gitlab)
		g.url = u
		g.apiUrl = "https://gitlab.com/api/v4"
		return g, nil
	default:
		return nil, errors.New(fmt.Sprintf("Undefined provider type: %s", u.Host))
	}
}
