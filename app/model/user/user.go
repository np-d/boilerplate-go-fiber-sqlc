package user

import "github.com/gofiber/fiber/v2"

type Router interface {
	Setup(app fiber.Router) fiber.Router
}

type Controller interface {
	Login(ctx *fiber.Ctx) error
	Get(ctx *fiber.Ctx) error
	Me(ctx *fiber.Ctx) error
	Create(ctx *fiber.Ctx) error
	Update(ctx *fiber.Ctx) error
	UpdatePassword(ctx *fiber.Ctx) error
	Delete(ctx *fiber.Ctx) error
}

type Service interface {
	Login(rp *LoginRequestPayload) (*map[string]any, error)
	Get(id int) (*map[string]any, error)
	Me(loggedUserId int) (*map[string]any, error)
	Create(rp *CreateRequestPayload) (*map[string]any, error)
	Update(loggedUserId int, rp *UpdateRequestPayload) (*map[string]any, error)
	UpdatePassword(loggedUserId int, rp *UpdatePasswordRequestPayload) error
	Delete(loggedUserId int, rp *DeleteRequestPayload) error
}

type LoginRequestPayload struct {
	Username string `json:"username" validate:"required,min=5,max=30"`
	Password string `json:"password" validate:"required,min=8"`
}

type CreateRequestPayload struct {
	DisplayName string `json:"display_name" validate:"required,min=1,max=50"`
	Username    string `json:"username" validate:"required,min=5,max=30"`
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=8"`
}

type UpdateRequestPayload struct {
	DisplayName string `json:"display_name" validate:"required,min=1,max=50"`
	Username    string `json:"username" validate:"required,min=5,max=30"`
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=8"`
}

type UpdatePasswordRequestPayload struct {
	OldPassword             string `json:"old_password" validate:"required,min=8"`
	NewPassword             string `json:"new_password" validate:"required,min=8"`
	NewPasswordConfirmation string `json:"new_password_confirmation" validate:"required,min=8"`
}

type DeleteRequestPayload struct {
	Password string `json:"password" validate:"required,min=8"`
}
