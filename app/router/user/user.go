package user

import (
	"github.com/gofiber/fiber/v2"
	controller "github.com/np-d/boilerplate-go-fiber-sqlc/app/controller/user"
	"github.com/np-d/boilerplate-go-fiber-sqlc/app/database/postgres"
	"github.com/np-d/boilerplate-go-fiber-sqlc/app/middleware/authorization"
	model "github.com/np-d/boilerplate-go-fiber-sqlc/app/model/user"
)

type RouterImpl struct {
	db *postgres.Database
}

func (r RouterImpl) Setup(app fiber.Router) fiber.Router {
	users := app.Group("/users")
	users.Post("/", controller.New(r.db).Create)
	users.Post("/login", controller.New(r.db).Login)

	users.Use(authorization.Authorize)

	users.Get("/me", controller.New(r.db).Me)
	users.Get("/:id", controller.New(r.db).Get)
	users.Delete("/", controller.New(r.db).Delete)
	users.Put("/", controller.New(r.db).Update)
	users.Patch("/password", controller.New(r.db).UpdatePassword)
	return users
}

func New(db *postgres.Database) model.Router {
	return &RouterImpl{db}
}
