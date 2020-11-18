package router

import (
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"github.com/nozgurozturk/marvin/server/internal/api"
	"github.com/nozgurozturk/marvin/server/internal/service"
	"github.com/nozgurozturk/marvin/server/internal/storage"
)

type Server struct {
	Router  *fiber.App
	Service service.Service
	Store   storage.Store
}

func New(s storage.Store) *Server {

	engine := html.New("./web", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})
	//app.Use(cors.New(cors.Config{
	//	AllowCredentials: true,
	//	AllowMethods:     "POST, OPTIONS, GET, PUT, DELETE, HEAD",
	//	AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
	//}))
	router := &Server{
		Router:  app,
		Service: service.New(s),
		Store:   s,
	}
	router.initializeRouters()
	return router
}

func (s *Server) initializeRouters() {

	authRouter := s.Router.Group("/auth")
	api.AuthHandler(authRouter, s.Service.Auth(), s.Service.User(), s.Service.Subscriber())

	apiRouter := s.Router.Group("/api", AuthMiddleware(s.Service.Auth()))

	userRouter := apiRouter.Group("/user")
	api.UserHandler(userRouter, s.Service.User(), s.Service.Repo())

	repoRouter := apiRouter.Group("/repository")
	api.RepositoryHandler(repoRouter, s.Service.Repo(), s.Service.Subscriber())

	subscriberRouter := apiRouter.Group("/subscriber")
	api.SubscriberHandler(subscriberRouter, s.Service.Subscriber(), s.Service.Repo())

	publicSubscriberRouter := s.Router.Group("/subscriber")
	api.PublicSubscriberHandler(publicSubscriberRouter, s.Service.Subscriber(), s.Service.Repo())

	// Documentation
	s.Router.Get("/docs/*", swagger.Handler)
}
