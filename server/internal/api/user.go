package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nozgurozturk/marvin/pkg/errors"
	"github.com/nozgurozturk/marvin/server/entity"
	"github.com/nozgurozturk/marvin/server/internal/app"
	"github.com/nozgurozturk/marvin/server/internal/service"
	"net/http"
)

func UserHandler(router fiber.Router, userService service.UserService, repoService service.RepoService) {
	router.Put("/", updateUser(userService))
	router.Delete("/", deleteUser(userService, repoService))
}

// updateUser is a function to update user values
// @Summary Updates user values (Accept partial updates)
// @Tags user
// @Accept json
// @Produce json
// @Param request body entity.UserDTO true "User"
// @Success 200 {object} entity.Response{data=entity.UserDTO}
// @Failure 401 {object} errors.AppError{}
// @Failure 422 {object} errors.AppError{}
// @Failure 500 {object} errors.AppError{}
// @Router /api/user [put]
func updateUser(s service.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var requestBody entity.UserDTO
		if err := c.BodyParser(&requestBody); err != nil {
			parseErr := errors.UnprocessableEntity("Invalid user body")
			return c.Status(parseErr.Status).JSON(parseErr)
		}

		updated, err := s.Update(&requestBody)
		if err != nil {
			return c.Status(err.Status).JSON(err)
		}

		response := entity.ToResponse("Successfully update your user info", http.StatusOK, updated)
		return c.Status(response.Status).JSON(response)
	}
}

// deleteUser is a function to remove user from store
// @Summary Removes user
// @Tags user
// @Produce json
// @Success 200 {object} entity.Response{}
// @Success 200 {object} entity.Response{}
// @Failure 401 {object} errors.AppError{}
// @Failure 500 {object} errors.AppError{}
// @Router /api/user [delete]
func deleteUser(s service.UserService, r service.RepoService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token, err := app.ExtractToken(c)
		if err != nil {
			return c.Status(err.Status).JSON(err)
		}

		claims, err := app.ExtractTokenMetaData(token)
		if err != nil {
			return c.Status(err.Status).JSON(err)
		}

		err = s.Delete(claims.UserID)
		if err != nil {
			return c.Status(err.Status).JSON(err)
		}

		err = r.DeleteMany(claims.UserID)
		if err != nil {
			return c.Status(err.Status).JSON(err)
		}

		response := entity.ToResponse("Deleted", http.StatusOK, nil)
		return c.Status(response.Status).JSON(response)
	}
}
