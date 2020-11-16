package providers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nozgurozturk/marvin/pkg/client"
	"net/url"
	"strings"
)

type Github struct {
	url    *url.URL
	apiUrl string
}

func (g *Github) UrlResolver() (string, string) {
	path := g.url.Path
	p := strings.Split(path, "/")
	return p[1], p[2]
}

func (g *Github) GetRepositoryTree(owner string, name string) ([]map[string]interface{}, error) {

	endpoint := fmt.Sprintf("/repos/%s/%s/contents", owner, name)
	headers := map[string]string{
		"Accept": "application/vnd.github.v3+json",
	}

	treeData, err := client.New(g.apiUrl).Get(endpoint, headers)
	if err != nil {
		return nil, err
	}

	if treeData == nil {
		return nil, errors.New("repository is not exist")
	}
	var tree []map[string]interface{}
	if err := json.Unmarshal(treeData, &tree); err != nil {
		return nil, err
	}

	return tree, nil
}


func (g *Github) FindPackagesInfo(tree []map[string]interface{}) []map[string]interface{} {

	var packagesInfo []map[string]interface{}

	for _, file := range tree {
		if file["name"] == "package.json" {
			packagesInfo = append(packagesInfo, file)
		}
		if file["name"] == "composer.json" {
			packagesInfo = append(packagesInfo, file)
		}
	}

	return packagesInfo
}

func (g *Github) GetPackageFiles(files []map[string]interface{}) (map[string]interface{}, error) {

	packageFiles := map[string]interface{}{}

	for _, file := range files {
		endpoint := file["download_url"].(string)
		headers := map[string]string{
			"Content-Type": "application/json",
		}

		packagesData, err := client.New("").Get(endpoint, headers)
		if err != nil {
			return nil, err
		}

		var packageFile map[string]interface{}
		if err := json.Unmarshal(packagesData, &packageFile); err != nil {
			return nil, err
		}

		fileName := file["name"].(string)
		if fileName != "" {
			packageFiles[fileName] = packageFile
		}
	}

	return packageFiles, nil
}
