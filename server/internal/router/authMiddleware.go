package router

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/nozgurozturk/marvin/pkg/errors"
	"github.com/nozgurozturk/marvin/server/internal/app"
	"github.com/nozgurozturk/marvin/server/internal/service"
)

func AuthMiddleware(authService service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		tokenString, err := app.ExtractToken(c)

		token, err := app.ValidateToken(tokenString, "Access")
		if err != nil {
			if token != nil {
				claims := token.Claims.(jwt.MapClaims)
				uuid, _ := claims["uuid"].(string)
				err := authService.DeleteAuth(uuid)
				if err != nil {
					return c.Status(err.Status).JSON(err)
				}
			}
			return c.Status(err.Status).JSON(err)
		}

		claims, _ := token.Claims.(jwt.MapClaims)
		uuid, _ := claims["uuid"].(string)
		userID, _ := claims["userID"].(string)
		authorized, _ := claims["authorized"].(bool)

		authID, err := authService.FindAuth(uuid)
		if err != nil {
			return c.Status(err.Status).JSON(err)
		}

		if authID != userID {
			notExistTokenErr := errors.Unauthorized("Unauthorized user")
			return c.Status(notExistTokenErr.Status).JSON(notExistTokenErr)
		}

		if !authorized {
			notAuthErr := errors.Unauthorized("Email confirmation is required")
			return c.Status(notAuthErr.Status).JSON(notAuthErr)
		}
		// Pass user id to handlers
		c.Locals("user", userID)
		return c.Next()
	}
}
