package entity

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type PackageVersion struct {
	Current string `json:"current" bson:"current"`
	Last    string `json:"last" bson:"last"`
}

type Package struct {
	Name       string         `json:"name" bson:"name"`
	Version    PackageVersion `json:"version" bson:"version"`
	File       string         `json:"file" bson:"file"`
	IsOutdated bool           `json:"isOutdated" bson:"isOutdated"`
}

type Repo struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserID      primitive.ObjectID `json:"userID" bson:"userID"`
	Name        string             `json:"name" bson:"name"`
	Owner       string             `json:"owner" bson:"owner"`
	Path        string             `json:"path" bson:"path"`
	Provider    string             `json:"provider" bson:"provider"`
	PackageList []*Package         `json:"packageList, omitempty" bson:"packageList,omitempty"`
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt"`
}

type RepoDTO struct {
	ID          *string    `json:"id,omitempty"`
	UserID      string     `json:"userID"`
	Name        string     `json:"name"`
	Owner       string     `json:"owner"`
	Path        string     `json:"path"`
	Provider    string     `json:"provider"`
	PackageList []*Package `json:"packageList, omitempty"`
}

type RepoIDRequest struct {
	ID string `json:"id"`
}

type RepoUrlRequest struct {
	Url string `json:"url"`
}

func ToRepoDTOs(repos []*Repo) []*RepoDTO {

	repoDTOs := make([]*RepoDTO, len(repos))

	for i, item := range repos {
		repoDTOs[i] = ToRepoDTO(item)
	}

	return repoDTOs
}

func ToRepoDTO(repo *Repo) *RepoDTO {

	id := repo.ID.Hex()

	return &RepoDTO{
		ID:          &id,
		UserID:      repo.UserID.Hex(),
		Name:        repo.Name,
		Owner:       repo.Owner,
		Path:        repo.Path,
		Provider:    repo.Provider,
		PackageList: repo.PackageList,
	}
}

func ToRepo(repoDTO *RepoDTO) *Repo {

	userId, _ := primitive.ObjectIDFromHex(repoDTO.UserID)

	repo := &Repo{
		Name:        repoDTO.Name,
		UserID:      userId,
		Owner:       repoDTO.Owner,
		Path:        repoDTO.Path,
		Provider:    repoDTO.Provider,
		PackageList: repoDTO.PackageList,
	}

	if repoDTO.ID != nil {
		repoId, _ := primitive.ObjectIDFromHex(*repoDTO.ID)
		repo.ID = repoId
	} else {
		repo.ID = primitive.NilObjectID
	}

	return repo
}

func ToPackageDTOs(rawPackages map[string]string, file string) []*Package {

	var packageDTOs []*Package

	for key, value := range rawPackages {
		packageDTOs = append(packageDTOs, &Package{
			Name: key,
			Version: PackageVersion{
				Current: value,
				Last:    value,
			},
			File:       file,
			IsOutdated: false,
		})
	}

	return packageDTOs
}
