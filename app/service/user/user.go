package user

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/np-d/boilerplate-go-fiber-sqlc/app/database/postgres"
	"github.com/np-d/boilerplate-go-fiber-sqlc/app/database/postgres/sqlc"
	model "github.com/np-d/boilerplate-go-fiber-sqlc/app/model/user"
	"github.com/np-d/boilerplate-go-fiber-sqlc/app/util/converter"
	"golang.org/x/crypto/bcrypt"
)

type ServiceImpl struct {
	db *postgres.Database
}

func (s ServiceImpl) Me(loggedUserId int) (*map[string]any, error) {
	user, err := s.db.Queries.GetUser(*s.db.Ctx, int32(loggedUserId))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fiber.NewError(fiber.StatusNotFound, "user not found")
		}
		return nil, fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	response, err := converter.StructToMap(&user)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	delete(*response, "password")
	delete(*response, "deleted_at")

	return response, nil
}

func (s ServiceImpl) Login(rp *model.LoginRequestPayload) (*map[string]any, error) {
	user, err := s.db.Queries.GetUserByUsername(*s.db.Ctx, rp.Username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(rp.Password)) != nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "invalid password")
	}

	expirationDate := time.Now().Add(time.Hour * 24).UTC()

	tokenData := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"iss": os.Getenv("JWT_ISSUER"),
		"sub": strconv.Itoa(int(user.ID)),
		"exp": expirationDate.Unix(),
		"iat": time.Now().UTC().Unix(),
	})

	token, err := tokenData.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return nil, err
	}

	return &map[string]any{
		"token": token,
		"exp":   expirationDate.Unix(),
	}, nil
}

func (s ServiceImpl) Get(id int) (*map[string]any, error) {
	user, err := s.db.Queries.GetUser(*s.db.Ctx, int32(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fiber.NewError(fiber.StatusNotFound, "user not found")
		}
		return nil, fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	response, err := converter.StructToMap(&user)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	delete(*response, "password")
	delete(*response, "email")
	delete(*response, "deleted_at")
	delete(*response, "updated_at")
	delete(*response, "created_at")

	return response, nil
}

func (s ServiceImpl) Create(rp *model.CreateRequestPayload) (*map[string]any, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(rp.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user, err := s.db.Queries.CreateUser(*s.db.Ctx, &sqlc.CreateUserParams{
		DisplayName: rp.DisplayName,
		Username:    rp.Username,
		Email:       rp.Email,
		Password:    string(hashedPassword),
	})
	if err != nil {
		return nil, fiber.NewError(fiber.StatusConflict, err.Error())
	}

	response, err := converter.StructToMap(user)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	delete(*response, "password")
	delete(*response, "deleted_at")

	return response, nil
}

func (s ServiceImpl) Update(loggedUserId int, rp *model.UpdateRequestPayload) (*map[string]any, error) {
	user, err := s.db.Queries.GetUser(*s.db.Ctx, int32(loggedUserId))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fiber.NewError(fiber.StatusNotFound, "user not found")
		}
		return nil, fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(rp.Password)) != nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "invalid password")
	}

	err = s.db.Queries.UpdateUser(*s.db.Ctx, &sqlc.UpdateUserParams{
		DisplayName: rp.DisplayName,
		Username:    rp.Username,
		Email:       rp.Email,
		ID:          int32(loggedUserId),
	})
	if err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	user, err = s.db.Queries.GetUser(*s.db.Ctx, int32(loggedUserId))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fiber.NewError(fiber.StatusNotFound, "user not found")
		}
		return nil, fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	response, err := converter.StructToMap(&user)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	delete(*response, "password")
	delete(*response, "deleted_at")

	return response, nil
}

func (s ServiceImpl) UpdatePassword(loggedUserId int, rp *model.UpdatePasswordRequestPayload) error {
	if rp.NewPassword != rp.NewPasswordConfirmation {
		return fiber.NewError(fiber.StatusBadRequest, "invalid new password confirmation")
	}

	user, err := s.db.Queries.GetUser(*s.db.Ctx, int32(loggedUserId))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fiber.NewError(fiber.StatusNotFound, "user not found")
		}
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(rp.OldPassword)) != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid password")
	}

	hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(rp.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	err = s.db.Queries.UpdateUserPassword(*s.db.Ctx, &sqlc.UpdateUserPasswordParams{
		ID:       int32(loggedUserId),
		Password: string(hashedNewPassword),
	})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return nil
}

func (s ServiceImpl) Delete(loggedUserId int, rp *model.DeleteRequestPayload) error {
	user, err := s.db.Queries.GetUser(*s.db.Ctx, int32(loggedUserId))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fiber.NewError(fiber.StatusNotFound, "user not found")
		}
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(rp.Password)) != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid password")
	}

	err = s.db.Queries.DeleteUser(*s.db.Ctx, int32(loggedUserId))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return nil
}

func New(db *postgres.Database) model.Service {
	return &ServiceImpl{db}
}
