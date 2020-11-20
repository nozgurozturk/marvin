package api

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/nozgurozturk/marvin/pkg/errors"
	"github.com/nozgurozturk/marvin/server/entity"
	"github.com/nozgurozturk/marvin/server/internal/app"
	"github.com/nozgurozturk/marvin/server/internal/config"
	"github.com/nozgurozturk/marvin/server/internal/service"
	"net/http"
)

// AuthHandler
func AuthHandler(router fiber.Router, authService service.AuthService, userService service.UserService, subService service.SubscriberService) {
	router.Post("/login", login(authService, userService))
	router.Post("/signup", signUp(authService, userService))
	router.Get("/logout", logout(authService))
	router.Post("/refresh", refresh(authService, userService))
	router.Get("/confirm", confirmAccount(authService, userService))
}

// login is a function to authenticate user
// @Summary Login with user credentials
// @Tags auth
// @Accept json
// @Produce json
// @Param request body entity.Login true "Login"
// @Success 200 {object} entity.Response{data=entity.TokenResponse}
// @Failure 401 {object} errors.AppError{}
// @Failure 422 {object} errors.AppError{}
// @Failure 500 {object} errors.AppError{}
// @Router /auth/login [post]
func login(authService service.AuthService, userService service.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		authLogin := new(entity.Login)
		if err := c.BodyParser(&authLogin); err != nil {
			parseError := errors.UnprocessableEntity("Invalid request body")
			return c.Status(parseError.Status).JSON(parseError)
		}

		currentUser, err := userService.FindByEmail(authLogin.Email)
		if err != nil {
			return c.Status(err.Status).JSON(err)
		}

		// Compares hashed password
		if passwordErr := entity.VerifyPassword(currentUser.Password, authLogin.Password); passwordErr != nil {
			err := errors.Unauthorized("Password is not correct")
			return c.Status(err.Status).JSON(err)
		}

		tokens, err := app.CreateToken(currentUser)
		if err != nil {
			return c.Status(err.Status).JSON(err)
		}

		err = authService.CreateAuth(tokens.RefreshToken)
		if err != nil {
			return c.Status(err.Status).JSON(err)
		}

		tokenStrings := entity.TokenDetailsToResponse(tokens)
		user := entity.ToUserResponse(currentUser)

		response := entity.ToResponse(
			"Successful login",
			http.StatusOK,
			fiber.Map{
				"tokens": tokenStrings,
				"user":   user,
			},
		)

		return c.Status(response.Status).JSON(response)
	}
}

// signup is a function to create new user and send email for verify
// @Summary Create new user and send verification email
// @Tags auth
// @Accept json
// @Produce json
// @Param request body entity.SignUp true "Sign up"
// @Success 200 {object} entity.Response{data=entity.TokenResponse}
// @Failure 409 {object} errors.AppError{}
// @Failure 422 {object} errors.AppError{}
// @Failure 500 {object} errors.AppError{}
// @Router /auth/signup [post]
func signUp(authService service.AuthService, userService service.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		u := new(entity.SignUp)
		if err := c.BodyParser(&u); err != nil {
			parseError := errors.UnprocessableEntity("Invalid request body")
			return c.Status(parseError.Status).JSON(parseError)
		}

		exist, err := userService.FindByEmail(u.Email)
		if err != nil && err.Status != http.StatusNotFound {
			return c.Status(err.Status).JSON(err)
		}

		if exist != nil {
			existErr := errors.AlreadyExist("This email is taken by an another user")
			return c.Status(existErr.Status).JSON(existErr)
		}

		createdUser, err := userService.Create(entity.AuthToUserDTO(u))
		if err != nil {
			return c.Status(err.Status).JSON(err)
		}

		tokens, err := app.CreateToken(createdUser)
		if err != nil {
			return c.Status(err.Status).JSON(err)
		}

		err = authService.CreateAuth(tokens.RefreshToken)
		if err != nil {
			return c.Status(err.Status).JSON(err)
		}

		// Gets serve configuration
		cnf := config.Get().HTTP

		// Creates email template data
		templateData := struct {
			User string
			Link string
		}{
			User: createdUser.Name,
			Link: fmt.Sprintf("http://%s%s/auth/confirm?t=%s", cnf.Host, cnf.Port, tokens.RefreshToken.Token),
		}

		emailBody, parsErr := app.ParseHTMLTemplate("./web/email-signup-verify.html", templateData)
		if parsErr != nil {
			e := errors.InternalServer(parsErr.Error())
			return c.Status(e.Status).JSON(e)
		}

		emailErr := app.SendEmail(createdUser.Email, "Marvin | Please confirm your email", emailBody)
		if emailErr != nil {
			e := errors.InternalServer(emailErr.Error())
			return c.Status(e.Status).JSON(e)
		}

		tokenStrings := entity.TokenDetailsToResponse(tokens)
		user := entity.ToUserResponse(createdUser)

		response := entity.ToResponse(
			"You successfully create an account. Please confirm your account with link that is sent to your email",
			http.StatusCreated,
			fiber.Map{
				"tokens": tokenStrings,
				"user":   user,
			},
		)
		return c.Status(response.Status).JSON(response)
	}
}

// logout is a function to drop user's session and delete session id
// @Summary Drops user's session
// @Tags auth
// @Produce json
// @Success 200 {object} entity.Response{}
// @Failure 401 {object} errors.AppError{}
// @Failure 500 {object} errors.AppError{}
// @Router /auth/logout [get]
func logout(authService service.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		t, err := app.ExtractToken(c)
		if err != nil {
			return c.Status(err.Status).JSON(err)
		}

		token, err := app.ValidateToken(t, "Access")
		claims := token.Claims.(jwt.MapClaims)
		uuid, _ := claims["uuid"].(string)

		if err != nil {
			if token != nil {
				err := authService.DeleteAuth(uuid)
				if err != nil {
					return c.Status(err.Status).JSON(err)
				}
			}
			return c.Status(err.Status).JSON(err)
		}

		err = authService.DeleteAuth(uuid)

		if err != nil {
			return c.Status(err.Status).JSON(err)
		}

		response := entity.ToResponse(
			"Successful logout",
			http.StatusOK,
			nil,
		)
		return c.Status(response.Status).JSON(response)
	}
}

// refresh is a function to renew user's session and delete old session id
// @Summary Renew user's session
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} entity.Response{data=entity.TokenResponse}
// @Failure 401 {object} errors.AppError{}
// @Failure 500 {object} errors.AppError{}
// @Router /auth/refresh [post]
func refresh(authService service.AuthService, userService service.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var requestBody struct {
			RefreshToken string `json:"refreshToken"`
		}

		if err := c.BodyParser(&requestBody); err != nil {
			parseError := errors.UnprocessableEntity("Invalid request body")
			return c.Status(parseError.Status).JSON(parseError)
		}

		if requestBody.RefreshToken == "" {
			err := errors.BadRequest("Refresh token is required")
			return c.Status(err.Status).JSON(err)
		}

		token, err := app.ValidateToken(requestBody.RefreshToken, "Refresh")
		claims := token.Claims.(jwt.MapClaims)
		uuid, _ := claims["uuid"].(string)

		if err != nil {
			if token != nil {
				err := authService.DeleteAuth(uuid)
				if err != nil {
					return c.Status(err.Status).JSON(err)
				}
			}
			return c.Status(err.Status).JSON(err)
		}

		userID, err := authService.FindAuth(uuid)
		if err != nil {
			return c.Status(err.Status).JSON(err)
		}

		foundUser, err := userService.FindByID(userID)
		if err != nil {
			return c.Status(err.Status).JSON(err)
		}

		tokens, err := app.CreateToken(foundUser)
		if err != nil {
			return c.Status(err.Status).JSON(err)
		}

		err = authService.CreateAuth(tokens.RefreshToken)
		if err != nil {
			return c.Status(err.Status).JSON(err)
		}

		err = authService.DeleteAuth(uuid)
		if err != nil {
			return c.Status(err.Status).JSON(err)
		}

		tokenResponse := entity.TokenDetailsToResponse(tokens)
		response := entity.ToResponse(
			"Successfully refreshes token",
			http.StatusOK,
			tokenResponse,
		)

		return c.Status(response.Status).JSON(response)
	}
}

// confirmAccount is a function to verify user's account with link that is sent by email
// @Summary Verify user's account
// @Tags auth
// @Param t query string true "token"
// @Produce html
// @Router /auth/confirm [get]
func confirmAccount(authService service.AuthService, userService service.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		tokenString := c.Query("t")
		token, err := app.ValidateToken(tokenString, "Refresh")
		claims := token.Claims.(jwt.MapClaims)
		uuid, _ := claims["uuid"].(string)

		if err != nil {
			if token != nil {
				err := authService.DeleteAuth(uuid)
				if err != nil {
					return c.Status(err.Status).Render("page-error", fiber.Map{
						"ErrorStatus":  err.Status,
						"ErrorMessage": err.Message,
						"Error":        err.Error,
					})
				}
			}
			return c.Status(err.Status).Render("page-error", fiber.Map{
				"ErrorStatus":  err.Status,
				"ErrorMessage": err.Message,
				"Error":        err.Error,
			})
		}

		userID, err := authService.FindAuth(uuid)
		if err != nil {
			return c.Status(err.Status).Render("page-error", fiber.Map{
				"ErrorStatus":  err.Status,
				"ErrorMessage": err.Message,
				"Error":        err.Error,
			})
		}

		foundUser, err := userService.FindByID(userID)
		if err != nil {
			return c.Status(err.Status).Render("page-error", fiber.Map{
				"ErrorStatus":  err.Status,
				"ErrorMessage": err.Message,
				"Error":        err.Error,
			})
		}

		err = userService.Confirm(foundUser.ID)
		if err != nil {
			return c.Status(err.Status).Render("page-error", fiber.Map{
				"ErrorStatus":  err.Status,
				"ErrorMessage": err.Message,
				"Error":        err.Error,
			})
		}

		tokens, err := app.CreateToken(foundUser)
		if err != nil {
			return c.Status(err.Status).Render("page-error", fiber.Map{
				"ErrorStatus":  err.Status,
				"ErrorMessage": err.Message,
				"Error":        err.Error,
			})
		}

		err = authService.CreateAuth(tokens.RefreshToken)
		if err != nil {
			return c.Status(err.Status).Render("page-error", fiber.Map{
				"ErrorStatus":  err.Status,
				"ErrorMessage": err.Message,
				"Error":        err.Error,
			})
		}

		err = authService.DeleteAuth(uuid)
		if err != nil {
			return c.Status(err.Status).Render("page-error", fiber.Map{
				"ErrorStatus":  err.Status,
				"ErrorMessage": err.Message,
				"Error":        err.Error,
			})
		}

		return c.Status(http.StatusOK).Render("page-signup-result", fiber.Map{
			"User": foundUser.Name,
		})
	}
}
