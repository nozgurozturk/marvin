package providers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nozgurozturk/marvin/pkg/client"
	"net/url"
	"strconv"
	"strings"
)

type Gitlab struct {
	url    *url.URL
	apiUrl string
}

func (g *Gitlab) UrlResolver() (string, string) {
	path := g.url.Path
	p := strings.Split(path, "/")
	return p[1], p[2]
}

// Gets Repository ID for consume gitlab's API for next requests
func (g *Gitlab) getRepositoryID(namespace string, name string) (string, error) {

	var repoId float64

	endpoint := fmt.Sprintf("/projects?search=%s", name)
	headers := map[string]string{
		"Content-Type": "application/json",
	}

	repoData, err := client.New(g.apiUrl).Get(endpoint, headers)
	if err != nil {
		return "", err
	}

	var repoList []map[string]interface{}
	if err := json.Unmarshal(repoData, &repoList); err != nil {
		panic(err)
	}

	pathWithNamespace := fmt.Sprintf("%s/%s", namespace, name)
	for _, repo := range repoList {
		if repo["path_with_namespace"] == pathWithNamespace {
			repoId = repo["id"].(float64)
			break
		}
	}

	id := strconv.Itoa(int(repoId))
	return id, nil
}

func (g *Gitlab) GetRepositoryTree(namespace string, name string) ([]map[string]interface{}, error) {

	projectID, err := g.getRepositoryID(namespace, name)
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("/projects/%s/repository/tree", projectID)
	headers := map[string]string{
		"Content-Type": "application/json",
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

	// add tree item to project id for getting package files
	for _, file := range tree {
		file["projectId"] = projectID
	}

	return tree, nil
}

func (g *Gitlab) FindPackagesInfo(tree []map[string]interface{}) []map[string]interface{} {

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

func (g *Gitlab) GetPackageFiles(files []map[string]interface{}) (map[string]interface{}, error) {

	packageFiles := map[string]interface{}{}

	for _, file := range files {
		endpoint := fmt.Sprintf("/projects/%s/repository/blobs/%s/raw", file["projectId"].(string), file["id"].(string))
		headers := map[string]string{
			"Content-Type": "application/json",
		}

		packagesData, err := client.New(g.apiUrl).Get(endpoint, headers)
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
