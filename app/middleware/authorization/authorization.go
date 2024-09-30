package authorization

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func Authorize(ctx *fiber.Ctx) error {
	tokenStr := ctx.Get("Authorization")

	invalidTokenError := fiber.NewError(fiber.StatusUnauthorized, "invalid token")

	if tokenStr == "" {
		return invalidTokenError
	}

	if strings.Contains(tokenStr, "Bearer") {
		tokenStr = strings.Trim(strings.Split(tokenStr, "Bearer")[1], " ")
	}

	tokenStr = strings.TrimSpace(tokenStr)

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != "HS512" {
			return nil, fiber.NewError(401, "unexpected signing method")
		}

		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	if !token.Valid {
		return invalidTokenError
	}

	userIdStr, err := token.Claims.GetSubject()
	if err != nil {
		return invalidTokenError
	}

	ctx.Locals("id", userIdStr)

	return ctx.Next()
}
