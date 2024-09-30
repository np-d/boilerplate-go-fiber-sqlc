package logger

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"log"
)

func Logger(ctx *fiber.Ctx) error {
	requestId := ctx.Locals("request_id")
	body := ctx.Body()
	if len(body) == 0 {
		log.Printf("%s %s %s %s", ctx.IP(), requestId, ctx.Method(), ctx.Path())
		return ctx.Next()
	}
	bodyJson := make(map[string]interface{})
	err := json.Unmarshal(body, &bodyJson)
	if err != nil {
		log.Println(err)
	}
	passwordFields := []string{"password", "new_password", "new_password_confirmation", "old_password"}
	for _, field := range passwordFields {
		if _, ok := bodyJson[field]; ok {
			delete(bodyJson, field)
		}
	}
	newBody, err := json.Marshal(bodyJson)
	if err != nil {
		log.Println(err)
	}
	log.Printf("%s %s %s %s %s", ctx.IP(), requestId, ctx.Method(), ctx.Path(), newBody)
	return ctx.Next()
}
