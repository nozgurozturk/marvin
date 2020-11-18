package service

import (
	"github.com/nozgurozturk/marvin/pkg/errors"
	"github.com/nozgurozturk/marvin/pkg/managers"
	"github.com/nozgurozturk/marvin/pkg/parsers"
	"github.com/nozgurozturk/marvin/pkg/providers"
	"github.com/nozgurozturk/marvin/pkg/utils"
	"github.com/nozgurozturk/marvin/server/entity"
	"github.com/nozgurozturk/marvin/server/internal/storage"
	"net/url"
	"sync"
)

// RepoService interface
type RepoService interface {
	// Create creates new git repository and saves into store
	Create(rawUrl string, userID string) (*entity.RepoDTO, *errors.AppError)
	// FindByID returns git repository with matching id
	FindByID(repoID string) (*entity.RepoDTO, *errors.AppError)
	// FindByUrlAndUserID returns git repository with matching url and user id
	FindByUrlAndUserID(url string, userID string) (*entity.RepoDTO, *errors.AppError)
	// FindAll returns git repository belongs to user
	FindAll(userID string) ([]*entity.RepoDTO, *errors.AppError)
	// UpdatePackages insert updated packages
	UpdatePackages(repoDTO *entity.RepoDTO) (*entity.RepoDTO, *errors.AppError)
	// Delete removes git repository
	Delete(repoID string) *errors.AppError
	// DeleteMany removes all git repositories belongs to user
	DeleteMany(userID string) *errors.AppError
}

type repoService struct {
	repository storage.RepoRepository
}

func NewRepoService(r storage.RepoRepository) RepoService {
	return &repoService{
		repository: r,
	}
}

/*
	1. Resolve Url -> owner, name
	2. Get Provider -> github
	3. Get ManagerType
	3. Get PackageManager -> npm
	4. Get Packages -> packages
	5. Get Each Package Version
	6. Compare Versions
	7. Create repo
*/
func (s *repoService) Create(rawUrl string, userID string) (*entity.RepoDTO, *errors.AppError) {

	// Parses rawUrl to url.URL
	u, err := url.Parse(rawUrl)
	if err != nil {
		return nil, errors.InternalServer(err.Error())
	}

	// Gets git provider with matching host name
	p, err := providers.GetProvider(u)
	if err != nil {
		return nil, errors.InternalServer(err.Error())
	}

	// Resolves git repository's owner and name from url
	owner, name := p.UrlResolver()

	// Gets git repository file tree from root directory
	tree, err := p.GetRepositoryTree(owner, name)
	if err != nil {
		return nil, errors.InternalServer(err.Error())
	}

	// Find package file info from provider's API
	packagesInfo := p.FindPackagesInfo(tree)
	if packagesInfo == nil {
		return nil, errors.InternalServer("Package file is not found")
	}

	// Gets packages from package file
	packageFiles, err := p.GetPackageFiles(packagesInfo)
	if err != nil {
		return nil, errors.InternalServer(err.Error())
	}

	var packages []*entity.Package

	// If git repository more than one package file with matching file names
	for pkgName, file := range packageFiles {

		// Creates new parser with matching package file name
		parser, err := parsers.NewParser(pkgName)
		if err != nil {
			return nil, errors.InternalServer(err.Error())
		}

		// Parses registries with name and versions
		rawPackages := parser.Parse(file.(map[string]interface{}))

		// Maps raw package array to entity.Package array
		pkgs := entity.ToPackageDTOs(rawPackages, pkgName)

		packages = append(packages, pkgs...)
	}

	var wg sync.WaitGroup
	for _, pkg := range packages {
		wg.Add(1)
		go func(pkg *entity.Package) {
			defer wg.Done()

			// Creates new package manager that is consume api
			m, err := managers.NewManager(pkg.File)
			if err != nil {
				return
			}

			// Gets latest registry version
			registryVersion, err := m.GetRegistryVersion(pkg.Name)
			if err != nil {
				return
			}

			// Compares semantic version of latest and current version
			isOutdated := utils.CompareVersions(registryVersion, pkg.Version.Current)

			if isOutdated {
				pkg.Version.Last = registryVersion
				pkg.IsOutdated = isOutdated
			}
		}(pkg)
	}
	wg.Wait()

	repo := &entity.RepoDTO{
		Name:        name,
		Owner:       owner,
		Path:        rawUrl,
		Provider:    u.Host,
		PackageList: packages,
		UserID:      userID,
	}

	createdRepo, err := s.repository.Create(entity.ToRepo(repo))
	if err != nil {
		return nil, errors.InternalServer(err.Error())
	}

	createdRepoDTO := entity.ToRepoDTO(createdRepo)

	return createdRepoDTO, nil

}

func (s *repoService) FindByID(repoID string) (*entity.RepoDTO, *errors.AppError) {

	repo, err := s.repository.FindByID(repoID)
	if err != nil {
		return nil, errors.NotFound("Repository is not found")
	}

	if repo == nil {
		return nil, errors.NotFound("Repository is not found")
	}

	return entity.ToRepoDTO(repo), nil
}

func (s *repoService) FindByUrlAndUserID(url string, userID string) (*entity.RepoDTO, *errors.AppError) {

	repo, err := s.repository.FindByUrlAndUserID(url, userID)
	if err != nil {
		return nil, errors.NotFound("Repository is not found")
	}

	if repo == nil {
		return nil, errors.NotFound("Repository is not found")
	}

	return entity.ToRepoDTO(repo), nil
}

func (s *repoService) FindAll(userID string) ([]*entity.RepoDTO, *errors.AppError) {

	repos, err := s.repository.FindAll(userID)

	if err != nil {
		return nil, errors.InternalServer(err.Error())
	}

	return entity.ToRepoDTOs(repos), nil
}

func (s *repoService) UpdatePackages(repoDTO *entity.RepoDTO) (*entity.RepoDTO, *errors.AppError) {

	// Parses rawUrl to url.URL
	u, err := url.Parse(repoDTO.Path)
	if err != nil {
		return nil, errors.InternalServer(err.Error())
	}

	// Gets git provider with matching host name
	p, err := providers.GetProvider(u)
	if err != nil {
		return nil, errors.InternalServer(err.Error())
	}

	// Resolves git repository's owner and name from url
	owner, name := p.UrlResolver()

	// Gets git repository file tree from root directory
	tree, err := p.GetRepositoryTree(owner, name)
	if err != nil {
		return nil, errors.InternalServer(err.Error())
	}

	// Find package file info from provider's API
	packagesInfo := p.FindPackagesInfo(tree)
	if packagesInfo == nil {
		return nil, errors.InternalServer("Package file is not found")
	}

	// Gets packages from package file
	packageFiles, err := p.GetPackageFiles(packagesInfo)
	if err != nil {
		return nil, errors.InternalServer(err.Error())
	}

	var packages []*entity.Package

	// If git repository more than one package file with matching file names
	for pkgName, file := range packageFiles {

		// Creates new parser with matching package file name
		parser, err := parsers.NewParser(pkgName)
		if err != nil {
			return nil, errors.InternalServer(err.Error())
		}

		// Parses registries with name and versions
		rawPackages := parser.Parse(file.(map[string]interface{}))

		// Maps raw package array to entity.Package array
		pkgs := entity.ToPackageDTOs(rawPackages, pkgName)

		packages = append(packages, pkgs...)
	}

	var wg sync.WaitGroup
	for _, pkg := range packages {
		wg.Add(1)
		go func(pkg *entity.Package) {
			defer wg.Done()

			// Creates new package manager that is consume api
			m, err := managers.NewManager(pkg.File)
			if err != nil {
				return
			}

			// Gets latest registry version
			registryVersion, err := m.GetRegistryVersion(pkg.Name)
			if err != nil {
				return
			}

			// Compares semantic version of latest and current version
			isOutdated := utils.CompareVersions(registryVersion, pkg.Version.Current)

			if isOutdated {
				pkg.Version.Last = registryVersion
				pkg.IsOutdated = isOutdated
			}
		}(pkg)
	}
	wg.Wait()

	repoDTO.PackageList = packages
	updated, err := s.repository.UpdatePackages(entity.ToRepo(repoDTO))
	if err != nil {
		return nil, errors.InternalServer(err.Error())
	}

	return entity.ToRepoDTO(updated), nil
}

func (s *repoService) Delete(repoID string) *errors.AppError {

	err := s.repository.Delete(repoID)

	if err != nil {
		return errors.InternalServer(err.Error())
	}

	return nil
}

func (s *repoService) DeleteMany(userID string) *errors.AppError {

	err := s.repository.DeleteMany(userID)

	if err != nil {
		return errors.InternalServer(err.Error())
	}

	return nil
}
