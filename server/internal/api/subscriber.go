package api

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/nozgurozturk/marvin/pkg/errors"
	"github.com/nozgurozturk/marvin/pkg/utils"
	"github.com/nozgurozturk/marvin/server/entity"
	"github.com/nozgurozturk/marvin/server/internal/app"
	"github.com/nozgurozturk/marvin/server/internal/config"
	"github.com/nozgurozturk/marvin/server/internal/service"
	"net/http"
)

// SubscriberHandler
func SubscriberHandler(router fiber.Router, subService service.SubscriberService, repoService service.RepoService) {
	router.Post("/", createSubscriber(subService, repoService))
	router.Post("/all", findAllSubscriber(subService))
	router.Delete("/", deleteSubscriber(subService))
	router.Post("/send", sendConfirm(subService, repoService))
}

// PublicSubscriberHandler no need authentication
func PublicSubscriberHandler(router fiber.Router, subService service.SubscriberService, repoService service.RepoService) {
	router.Get("/confirm", confirmSub(subService, repoService))
	router.Get("/unsubscribe", unsubscribeSub(subService))
	router.Put("/update", updateSubscriberNotify(subService))
	router.Get("/", updateSubscriber(subService, repoService))
}

// createSubscriber is a function to create new subscriber
// @Summary Create subscriber belongs to repository and send email
// @Tags subscriber
// @Accept json
// @Produce json
// @Param request body entity.SubscriberRequest true "Subscriber"
// @Success 201 {object} entity.Response{data=entity.SubscriberDTO}
// @Failure 401 {object} errors.AppError{}
// @Failure 422 {object} errors.AppError{}
// @Failure 500 {object} errors.AppError{}
// @Router /api/subscriber [post]
func createSubscriber(s service.SubscriberService, r service.RepoService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		requestBody := new(entity.SubscriberRequest)
		if err := c.BodyParser(&requestBody); err != nil {
			e := errors.UnprocessableEntity("Invalid request body")
			return c.Status(e.Status).JSON(e)
		}

		repository, err := r.FindByID(*requestBody.RepoID)
		if err != nil {
			return c.Status(err.Status).JSON(err)
		}

		if c.Locals("user") != repository.UserID {
			err = errors.Forbidden("You don't have access")
			return c.Status(err.Status).JSON(err)
		}

		exist, err := s.FindByEmailAndRepoID(*requestBody.Email, *requestBody.RepoID)
		if err != nil && err.Status != http.StatusNotFound {
			return c.Status(err.Status).JSON(err)
		}

		if exist != nil {
			existErr := errors.AlreadyExist("Subscriber is already exist")
			return c.Status(existErr.Status).JSON(existErr)
		}

		subscriber, err := s.Create(*requestBody.Email, *requestBody.RepoID)
		if err != nil {
			return c.Status(err.Status).JSON(err)
		}

		token, err := app.CreateSubToken(subscriber)
		if err != nil {
			return c.Status(err.Status).JSON(err)
		}

		cnf := config.Get().HTTP

		templateData := struct {
			RepoName string
			RepoLink string
			Link     string
		}{
			RepoName: repository.Name,
			RepoLink: repository.Path,
			Link:     "http://" + cnf.Host + cnf.Port + "/subscriber/confirm?t=" + token.Token,
		}

		emailBody, parsErr := utils.ParseHTMLTemplate("./web/email-sub-confirm.html", templateData)
		if parsErr != nil {
			e := errors.InternalServer(parsErr.Error())
			return c.Status(e.Status).JSON(e)
		}

		emailErr := app.SendEmail(subscriber.Email, "Marvin Subscription Confirm", emailBody)
		if emailErr != nil {
			e := errors.InternalServer(emailErr.Error())
			return c.Status(e.Status).JSON(e)
		}

		response := entity.ToResponse(
			"You successfully add a new subscriber.",
			http.StatusCreated,
			subscriber,
		)

		return c.Status(response.Status).JSON(response)
	}
}

// findAllSubscriber is a function to find all subscriber
// @Summary Returns all subscriber belongs to repository
// @Tags subscriber
// @Accept json
// @Produce json
// @Param request body entity.RepoIDRequest true "Id"
// @Success 200 {object} entity.Response{data=[]entity.SubscriberDTO}
// @Failure 401 {object} errors.AppError{}
// @Failure 422 {object} errors.AppError{}
// @Failure 500 {object} errors.AppError{}
// @Router /api/subscriber/all [post]
func findAllSubscriber(s service.SubscriberService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		requestBody := new(entity.RepoIDRequest)

		if err := c.BodyParser(&requestBody); err != nil {
			e := errors.UnprocessableEntity("Invalid request body")
			return c.Status(e.Status).JSON(e)
		}

		found, err := s.FindAll(requestBody.ID)

		if err != nil {
			return c.Status(err.Status).JSON(err)
		}

		response := entity.ToResponse(
			"You successfully add a new subscriber.",
			http.StatusCreated,
			found,
		)
		return c.Status(response.Status).JSON(response)
	}
}

// deleteSubscriber is a function to remove subscriber
// @Summary Remove subscriber from repository
// @Tags subscriber
// @Accept json
// @Produce json
// @Param request body entity.SubscriberIDRequest true "Id"
// @Success 200 {object} entity.Response{}
// @Failure 401 {object} errors.AppError{}
// @Failure 422 {object} errors.AppError{}
// @Failure 500 {object} errors.AppError{}
// @Router /api/subscriber [delete]
func deleteSubscriber(s service.SubscriberService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		requestBody := new(entity.SubscriberIDRequest)
		if err := c.BodyParser(&requestBody); err != nil {
			e := errors.UnprocessableEntity("Invalid request body")
			return c.Status(e.Status).JSON(e)
		}

		err := s.Delete(requestBody.ID)
		if err != nil {
			return c.Status(err.Status).JSON(err)
		}

		response := entity.ToResponse("Subscriber has been deleted", http.StatusOK, nil)
		return c.Status(response.Status).JSON(response)
	}
}

// sendConfirm is a function to send confirm email to existing subscriber
// @Summary Send confirm email to existing subscriber
// @Tags subscriber
// @Accept json
// @Produce json
// @Param request body entity.SubscriberRequest true "Subscriber"
// @Success 201 {object} entity.Response{data=entity.SubscriberDTO}
// @Failure 401 {object} errors.AppError{}
// @Failure 422 {object} errors.AppError{}
// @Failure 500 {object} errors.AppError{}
// @Router /api/subscriber/send [post]
func sendConfirm(s service.SubscriberService, r service.RepoService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		requestBody := new(entity.SubscriberRequest)
		if err := c.BodyParser(&requestBody); err != nil {
			e := errors.UnprocessableEntity("Invalid request body")
			return c.Status(e.Status).JSON(e)
		}

		repository, err := r.FindByID(*requestBody.RepoID)
		if err != nil {
			return c.Status(err.Status).JSON(err)
		}

		if c.Locals("user") != repository.UserID {
			err = errors.Forbidden("You don't have access")
			return c.Status(err.Status).JSON(err)
		}

		exist, err := s.FindByEmailAndRepoID(*requestBody.Email, *requestBody.RepoID)
		if err != nil {
			return c.Status(err.Status).JSON(err)
		}

		token, err := app.CreateSubToken(exist)
		if err != nil {
			return c.Status(err.Status).JSON(err)
		}

		cnf := config.Get().HTTP

		templateData := struct {
			RepoName string
			RepoLink string
			Link     string
		}{
			RepoName: repository.Name,
			RepoLink: repository.Path,
			Link:     "http://" + cnf.Host + cnf.Port + "/subscriber/confirm?t=" + token.Token,
		}

		emailBody, parsErr := utils.ParseHTMLTemplate("./web/email-sub-confirm.html", templateData)
		if parsErr != nil {
			e := errors.InternalServer(parsErr.Error())
			return c.Status(e.Status).JSON(e)
		}

		emailErr := app.SendEmail(exist.Email, "Marvin Subscription Confirm", emailBody)
		if emailErr != nil {
			e := errors.InternalServer(emailErr.Error())
			return c.Status(e.Status).JSON(e)
		}

		response := entity.ToResponse(
			"You successfully sent an email to subscriber.",
			http.StatusCreated,
			exist,
		)

		return c.Status(response.Status).JSON(response)
	}
}

// confirmSub is a function to verify subscriber's email
// @Summary Verify subscriber's email
// @Tags subscriber
// @Produce html
// @Param t query string true "token"
// @Success 200 {object} entity.Response{}
// @Failure 403 {object} errors.AppError{}
// @Failure 422 {object} errors.AppError{}
// @Failure 500 {object} errors.AppError{}
// @Router /subscriber/confirm [get]
func confirmSub(s service.SubscriberService, r service.RepoService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := c.Query("t")
		token, err := app.ValidateToken(tokenString, "Sub")
		claims := token.Claims.(jwt.MapClaims)
		id, _ := claims["subID"].(string)
		email, _ := claims["email"].(string)
		repoID, _ := claims["repoID"].(string)

		subscriber, err := s.FindByID(id)
		if err != nil {
			return c.Status(err.Status).Render("page-error", fiber.Map{
				"ErrorStatus":  err.Status,
				"ErrorMessage": err.Message,
				"Error":        err.Error,
			})
		}

		repository, err := r.FindByID(repoID)
		if err != nil {
			return c.Status(err.Status).Render("page-error", fiber.Map{
				"ErrorStatus":  err.Status,
				"ErrorMessage": err.Message,
				"Error":        err.Error,
			})
		}

		if subscriber.Email != email {
			err = errors.Forbidden("You don't have access")
			return c.Status(err.Status).Render("page-error", fiber.Map{
				"ErrorStatus":  err.Status,
				"ErrorMessage": err.Message,
				"Error":        err.Error,
			})
		}

		if subscriber.RepoID != repoID {
			err = errors.Forbidden("You don't have access")
			return c.Status(err.Status).Render("page-error", fiber.Map{
				"ErrorStatus":  err.Status,
				"ErrorMessage": err.Message,
				"Error":        err.Error,
			})
		}
		notifyTime := fmt.Sprintf("%02d:%02d", subscriber.Notify.Hour, subscriber.Notify.Minute)

		err = s.Confirm(subscriber, true)
		if err != nil {
			return c.Status(err.Status).JSON(err)
		}

		return c.Status(http.StatusOK).Render("page-sub-result", fiber.Map{
			"Email":      subscriber.Email,
			"NotifyTime": notifyTime,
			"RepoLink":   repository.Path,
			"RepoName":   repository.Name,
			"Hour":       subscriber.Notify.Hour,
			"Minute":     subscriber.Notify.Minute,
			"Weekday":    subscriber.Notify.Weekday,
			"Frequency":  subscriber.Notify.Frequency,
		})
	}
}

// unsubscribeSub is a function to unsubscribe subscriber's email from list
// @Summary Unsubscribe subscriber's email
// @Tags subscriber
// @Produce html
// @Param t query string true "token"
// @Success 200 {object} entity.Response{}
// @Failure 403 {object} errors.AppError{}
// @Failure 422 {object} errors.AppError{}
// @Failure 500 {object} errors.AppError{}
// @Router /subscriber/unsubscribe [get]
func unsubscribeSub(s service.SubscriberService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := c.Query("t")
		token, err := app.ValidateToken(tokenString, "Sub")
		claims := token.Claims.(jwt.MapClaims)
		id, _ := claims["subID"].(string)
		email, _ := claims["email"].(string)
		repoID, _ := claims["repoID"].(string)

		subscriber, err := s.FindByID(id)
		if err != nil {
			return c.Status(err.Status).Render("page-error", fiber.Map{
				"ErrorStatus":  err.Status,
				"ErrorMessage": err.Message,
				"Error":        err.Error,
			})
		}

		if subscriber.Email != email {
			err = errors.Forbidden("You don't have access")
			return c.Status(err.Status).Render("page-error", fiber.Map{
				"ErrorStatus":  err.Status,
				"ErrorMessage": err.Message,
				"Error":        err.Error,
			})
		}

		if subscriber.RepoID != repoID {
			err = errors.Forbidden("You don't have access")
			return c.Status(err.Status).Render("page-error", fiber.Map{
				"ErrorStatus":  err.Status,
				"ErrorMessage": err.Message,
				"Error":        err.Error,
			})
		}

		err = s.Confirm(subscriber, false)
		if err != nil {
			return c.Status(err.Status).JSON(err)
		}

		return c.Status(http.StatusOK).Render("page-sub-unsubscribe", fiber.Map{})
	}
}

// updateSubscriber is a function to render subscriber's notification preferences page
// @Summary Render subscriber's notification preferences page
// @Tags subscriber
// @Produce html
// @Param t query string true "token"
// @Success 200 {object} entity.Response{}
// @Failure 403 {object} errors.AppError{}
// @Failure 422 {object} errors.AppError{}
// @Failure 500 {object} errors.AppError{}
// @Router /subscriber [get]
func updateSubscriber(s service.SubscriberService, r service.RepoService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := c.Query("t")
		token, err := app.ValidateToken(tokenString, "Sub")
		claims := token.Claims.(jwt.MapClaims)
		id, _ := claims["subID"].(string)
		email, _ := claims["email"].(string)
		repoID, _ := claims["repoID"].(string)

		subscriber, err := s.FindByID(id)
		if err != nil {
			return c.Status(err.Status).JSON(err)
		}

		repository, err := r.FindByID(repoID)
		if err != nil {
			return c.Status(err.Status).JSON(err)
		}

		if subscriber.Email != email {
			return c.SendStatus(http.StatusForbidden)
		}

		if subscriber.RepoID != repoID {
			return c.SendStatus(http.StatusForbidden)
		}
		notifyTime := fmt.Sprintf("%02d:%02d", subscriber.Notify.Hour, subscriber.Notify.Minute)

		return c.Status(http.StatusOK).Render("page-sub-update", fiber.Map{
			"Email":      subscriber.Email,
			"NotifyTime": notifyTime,
			"RepoLink":   repository.Path,
			"RepoName":   repository.Name,
			"Hour":       subscriber.Notify.Hour,
			"Minute":     subscriber.Notify.Minute,
			"Weekday":    subscriber.Notify.Weekday,
			"Frequency":  subscriber.Notify.Frequency,
		})
	}
}

// updateSubscriberNotify is a function to update subscriber's notification preferences
// @Summary Update subscriber's notification preferences
// @Tags subscriber
// @Accept json
// @Produce json
// @Param request body entity.SubscriberRequest true "Subscriber"
// @Success 200
// @Failure 403 {object} errors.AppError{}
// @Failure 422 {object} errors.AppError{}
// @Failure 500 {object} errors.AppError{}
// @Router /subscriber/update [put]
func updateSubscriberNotify(s service.SubscriberService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		requestBody := new(entity.SubscriberRequest)
		if err := c.BodyParser(&requestBody); err != nil {
			e := errors.UnprocessableEntity("Invalid request body")
			return c.Status(e.Status).JSON(e)
		}

		tokenString := c.Query("t")
		_, err := app.ValidateToken(tokenString, "Sub")
		if err != nil {
			return c.Status(err.Status).JSON(err)
		}

		claims, err := app.ExtractSubTokenMetaData(tokenString)
		if err != nil {
			return c.Status(err.Status).JSON(err)
		}

		subscriber, err := s.FindByID(claims.SubID)
		if err != nil {
			return c.Status(err.Status).JSON(err)
		}

		if subscriber.Email != claims.Email {
			return c.SendStatus(http.StatusForbidden)
		}

		if subscriber.RepoID != claims.RepoID {
			return c.SendStatus(http.StatusForbidden)
		}

		subscriber.Notify = requestBody.Notify
		_, err = s.Update(subscriber)
		if err != nil {
			return c.Status(err.Status).JSON(err)
		}

		return c.SendStatus(http.StatusOK)
	}
}
