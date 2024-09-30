package user

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/np-d/boilerplate-go-fiber-sqlc/app/database/postgres"
	model "github.com/np-d/boilerplate-go-fiber-sqlc/app/model/user"
	service "github.com/np-d/boilerplate-go-fiber-sqlc/app/service/user"
)

type ControllerImpl struct {
	db *postgres.Database
}

var validate = validator.New()

func (c ControllerImpl) Me(ctx *fiber.Ctx) error {
	loggedUserId, _ := strconv.Atoi(ctx.Locals("id").(string))
	resp, err := service.New(c.db).Me(loggedUserId)
	if err != nil {
		return err
	}
	return ctx.JSON(resp)
}

func (c ControllerImpl) Login(ctx *fiber.Ctx) error {
	req := model.LoginRequestPayload{}
	err := json.Unmarshal(ctx.Body(), &req)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if errors.As(validate.Struct(&req), &err) {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	resp, err := service.New(c.db).Login(&req)
	if err != nil {
		return err
	}
	return ctx.JSON(resp)
}

func (c ControllerImpl) Get(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	resp, err := service.New(c.db).Get(id)
	if err != nil {
		return err
	}
	return ctx.JSON(resp)
}

func (c ControllerImpl) Create(ctx *fiber.Ctx) error {
	req := model.CreateRequestPayload{}
	err := json.Unmarshal(ctx.Body(), &req)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if errors.As(validate.Struct(&req), &err) {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	resp, err := service.New(c.db).Create(&req)
	if err != nil {
		return err
	}
	return ctx.JSON(resp)
}

func (c ControllerImpl) Update(ctx *fiber.Ctx) error {
	loggedUserId, _ := strconv.Atoi(ctx.Locals("id").(string))
	req := model.UpdateRequestPayload{}
	err := json.Unmarshal(ctx.Body(), &req)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	err = validate.Struct(&req)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	resp, err := service.New(c.db).Update(loggedUserId, &req)
	if err != nil {
		return err
	}
	return ctx.JSON(resp)
}

func (c ControllerImpl) UpdatePassword(ctx *fiber.Ctx) error {
	loggedUserId, _ := strconv.Atoi(ctx.Locals("id").(string))
	req := model.UpdatePasswordRequestPayload{}
	err := json.Unmarshal(ctx.Body(), &req)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if errors.As(validate.Struct(&req), &err) {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	err = service.New(c.db).UpdatePassword(loggedUserId, &req)
	if err != nil {
		return err
	}
	return ctx.JSON(fiber.Map{})
}

func (c ControllerImpl) Delete(ctx *fiber.Ctx) error {
	loggedUserId, _ := strconv.Atoi(ctx.Locals("id").(string))
	req := model.DeleteRequestPayload{}
	err := json.Unmarshal(ctx.Body(), &req)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if errors.As(validate.Struct(&req), &err) {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	err = service.New(c.db).Delete(loggedUserId, &req)
	if err != nil {
		return err
	}
	return ctx.JSON(fiber.Map{})
}

func New(db *postgres.Database) model.Controller {
	return &ControllerImpl{db}
}
