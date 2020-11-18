package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nozgurozturk/marvin/pkg/errors"
	"github.com/nozgurozturk/marvin/server/entity"
	"github.com/nozgurozturk/marvin/server/internal/app"
	"github.com/nozgurozturk/marvin/server/internal/service"
	"net/http"
)

func RepositoryHandler(router fiber.Router, repoService service.RepoService, subService service.SubscriberService) {
	router.Post("/", createRepo(repoService))
	router.Get("/", findAllRepo(repoService))
	router.Put("/", updateRepoPackages(repoService))
	router.Delete("/", deleteRepo(repoService, subService))
}

// createRepo is a function to create new git repository
// @Summary Create new git repository with packages
// @Tags repo
// @Accept json
// @Produce json
// @Param request body entity.RepoUrlRequest true "Url"
// @Success 201 {object} entity.Response{data=entity.RepoDTO}
// @Failure 401 {object} errors.AppError{}
// @Failure 422 {object} errors.AppError{}
// @Failure 500 {object} errors.AppError{}
// @Router /api/repository [post]
func createRepo(s service.RepoService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		requestBody := new(entity.RepoUrlRequest)
		if err := c.BodyParser(&requestBody); err != nil {
			e := errors.UnprocessableEntity("Invalid request body")
			return c.Status(e.Status).JSON(e)
		}

		token, err := app.ExtractToken(c)
		if err != nil {
			return c.Status(err.Status).JSON(err)
		}

		claims, err := app.ExtractTokenMetaData(token)
		if err != nil {
			return c.Status(err.Status).JSON(err)
		}

		exist, _ := s.FindByUrlAndUserID(requestBody.Url, claims.UserID)
		if exist != nil {
			err = errors.AlreadyExist("Repository is already exist")
			return c.Status(err.Status).JSON(err)
		}

		repo, err := s.Create(requestBody.Url, claims.UserID)
		if err != nil {
			return c.Status(err.Status).JSON(err)
		}

		response := entity.ToResponse(
			"You successfully create an git repository.",
			http.StatusCreated,
			repo,
		)
		return c.Status(response.Status).JSON(response)
	}
}

// findAllRepo is a function to returns all git repository that user have
// @Summary Returns all git repository that user have
// @Tags repo
// @Accept json
// @Produce json
// @Success 200 {object} entity.Response{data=[]entity.RepoDTO}
// @Failure 401 {object} errors.AppError{}
// @Failure 422 {object} errors.AppError{}
// @Failure 500 {object} errors.AppError{}
// @Router /api/repository [get]
func findAllRepo(s service.RepoService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		token, err := app.ExtractToken(c)
		if err != nil {

			return c.Status(err.Status).JSON(err)
		}

		claims, err := app.ExtractTokenMetaData(token)
		if err != nil {

			return c.Status(err.Status).JSON(err)
		}

		found, err := s.FindAll(claims.UserID)
		if err != nil {
			return c.Status(err.Status).JSON(err)
		}

		response := entity.ToResponse(
			"All repositories that you have",
			http.StatusOK,
			found,
		)
		return c.Status(response.Status).JSON(response)
	}
}

// updateRepoPackages is a function to update repository's dependencies
// @Summary Updates dependencies and compare versions
// @Tags repo
// @Accept json
// @Produce json
// @Param request body entity.RepoIDRequest true "Id"
// @Success 200 {object} entity.Response{data=entity.RepoDTO}
// @Failure 401 {object} errors.AppError{}
// @Failure 404 {object} errors.AppError{}
// @Failure 422 {object} errors.AppError{}
// @Failure 500 {object} errors.AppError{}
// @Router /api/repository [put]
func updateRepoPackages(s service.RepoService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		requestBody := new(entity.RepoIDRequest)

		if err := c.BodyParser(&requestBody); err != nil {
			e := errors.UnprocessableEntity("Invalid request body")
			return c.Status(e.Status).JSON(e)
		}

		repo, err := s.FindByID(requestBody.ID)
		if err != nil {
			return c.Status(err.Status).JSON(err)
		}

		if c.Locals("user") != repo.UserID {
			err = errors.Forbidden("You don't have access")
			return c.Status(err.Status).JSON(err)
		}

		updated, err := s.UpdatePackages(repo)
		if err != nil {
			return c.Status(err.Status).JSON(err)
		}

		response := entity.ToResponse(
			"Packages are updated.",
			http.StatusOK,
			updated,
		)

		return c.Status(response.Status).JSON(response)
	}
}

// deleteRepo is a function to remove repository from database
// @Summary Remove repository and subscribers belongs to it
// @Tags repo
// @Accept json
// @Produce json
// @Param request body entity.RepoIDRequest true "Id"
// @Success 200 {object} entity.Response{}
// @Failure 401 {object} errors.AppError{}
// @Failure 404 {object} errors.AppError{}
// @Failure 422 {object} errors.AppError{}
// @Failure 500 {object} errors.AppError{}
// @Router /api/repository [delete]
func deleteRepo(s service.RepoService, sub service.SubscriberService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		requestBody := new(entity.RepoIDRequest)

		if err := c.BodyParser(&requestBody); err != nil {
			e := errors.UnprocessableEntity("Invalid request body")
			return c.Status(e.Status).JSON(e)
		}

		repo, err := s.FindByID(requestBody.ID)
		if err != nil {
			return c.Status(err.Status).JSON(err)
		}

		if c.Locals("user") != repo.UserID {
			err = errors.Forbidden("You don't have access")
			return c.Status(err.Status).JSON(err)
		}

		err = s.Delete(requestBody.ID)
		if err != nil {
			return c.Status(err.Status).JSON(err)
		}

		err = sub.DeleteAll(requestBody.ID)
		if err != nil {
			return c.Status(err.Status).JSON(err)
		}

		response := entity.ToResponse(
			"Repository and all subscribers that repository have, has been deleted",
			http.StatusOK,
			nil,
		)
		return c.Status(response.Status).JSON(response)
	}
}
