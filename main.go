package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/google/uuid"
	"github.com/np-d/boilerplate-go-fiber-sqlc/app/middleware/logger"

	"github.com/np-d/boilerplate-go-fiber-sqlc/app/database/postgres"
	userRouter "github.com/np-d/boilerplate-go-fiber-sqlc/app/router/user"
)

func main() {
	if !checkEnvironmentVariables() {
		log.Fatalln("please check your environment variables")
	}

	db := connectToDatabase()
	defer db.Close()

	app := fiber.New(fiber.Config{
		Prefork:       true,
		CaseSensitive: true,
		ServerHeader:  os.Getenv("APP_HEADER"),
		AppName:       os.Getenv("APP_NAME"),
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError

			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}

			ctx.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

			requestId := ctx.Locals("request_id").(string)

			log.Printf("%s %s %s %s | %s", ctx.IP(), requestId, ctx.Method(), ctx.Path(), err.Error())

			return ctx.Status(code).JSON(fiber.Map{"request_id": requestId, "error": err.Error()})
		},
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowCredentials: false,
	}))

	app.Use(setRequestId)
	app.Use(logger.Logger)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(map[string]string{
			"message": "Welcome to your_project!",
		})
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(map[string]string{
			"status": "OK",
		})
	})

	api := app.Group("api")

	v1 := api.Group("v1")

	userRouter.New(db).Setup(v1)

	var err error
	switch os.Getenv("APP_ENV") {
	case "prod":
		err = app.ListenTLS(fmt.Sprintf("%s:%s", os.Getenv("APP_HOST"), os.Getenv("APP_PORT")), "./certificate/certificate.crt", "./certificate/private.key")
	case "dev":
		err = app.Listen(fmt.Sprintf("%s:%s", os.Getenv("APP_HOST"), os.Getenv("APP_PORT")))
	}

	if err != nil {
		panic(err)
	}
}

func setRequestId(ctx *fiber.Ctx) error {
	ctx.Locals("request_id", uuid.New().String())
	return ctx.Next()
}

func checkEnvironmentVariables() bool {
	switch {
	case os.Getenv("APP_ENV") == "":
		return false
	case os.Getenv("APP_HOST") == "":
		return false
	case os.Getenv("APP_PORT") == "":
		return false
	case os.Getenv("APP_HEADER") == "":
		return false
	case os.Getenv("APP_NAME") == "":
		return false
	case os.Getenv("DATABASE_URL") == "":
		return false
	case os.Getenv("JWT_SECRET") == "":
		return false
	case os.Getenv("JWT_ISSUER") == "":
		return false
	default:
		return true
	}
}

func connectToDatabase() *postgres.Database {
	return postgres.Connect()
}
